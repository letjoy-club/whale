package loader

import (
	"context"
	"sync"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/models"

	"github.com/letjoy-club/mida-tool/logger"
	"github.com/letjoy-club/mida-tool/shortid"
	"github.com/patrickmn/go-cache"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserDiscoverMatching struct {
	userID string

	viewedMatchingIDs map[string]struct{}

	db              *gorm.DB
	mu              sync.Mutex
	matchingMap     GroupedMatchings
	cityMatchingMap map[string]CityMatchings

	nextToken map[string]int
}

type UserDiscoverOpt struct {
	CityID string
	Gender models.Gender
	N      int

	NextToken string
}

func (u *UserDiscoverMatching) LoadMatchingIDs(
	ctx context.Context,
	matchingLoader func(topicID, cityID string) []*models.Matching,
	topicID string,
	opt UserDiscoverOpt,
) (retIDs []string, next string) {
	u.mu.Lock()
	defer u.mu.Unlock()

	// 如果查看的数量多于 300，就不再返回
	if len(u.viewedMatchingIDs) > 300 {
		return []string{}, ""
	}

	startIndex := u.nextToken[opt.NextToken]

	var matchings []*models.Matching
	var ok bool

	if opt.CityID == "" {
		matchings, ok = u.matchingMap[topicID]
		if !ok {
			matchings = matchingLoader(topicID, opt.CityID)
			u.matchingMap[topicID] = matchings
		}
	} else {
		cityMatchings, ok := u.cityMatchingMap[opt.CityID]
		if !ok {
			u.cityMatchingMap[opt.CityID] = CityMatchings{}
			cityMatchings = u.cityMatchingMap[opt.CityID]
		}

		if cityMatchings[topicID] == nil {
			matchings = matchingLoader(topicID, opt.CityID)
			u.cityMatchingMap[opt.CityID][topicID] = matchings
		}
	}

	retN := 0
	retIDs = make([]string, 0, opt.N)

	for i := startIndex; i < len(matchings); i++ {
		matching := matchings[i]
		// 过滤自己的
		if matching.UserID == u.userID {
			continue
		}
		// 城市过滤
		if opt.CityID == "" || opt.CityID == matching.CityID {
			// 性别过滤
			if opt.Gender != models.GenderN && matching.MyGender != opt.Gender.String() {
				continue
			}
			retN++
			retIDs = append(retIDs, matching.ID)
			if retN >= opt.N {
				next = shortid.New("", 4)
				u.nextToken[next] = i + 1
				break
			}
		}
	}
	// 记录已经看过的
	for _, id := range retIDs {
		u.viewedMatchingIDs[id] = struct{}{}
	}
	return retIDs, next
}

type GroupedMatchings map[string][]*models.Matching
type CityMatchings = GroupedMatchings

func (g GroupedMatchings) Load(topicID, gender string, start, n int) []string {
	retN := 0
	retIDs := make([]string, 0, n)

	if g[topicID] == nil {
		return []string{}
	}
	for _, matching := range g[topicID] {
		if gender == models.GenderN.String() || gender == matching.MyGender {
			retIDs = append(retIDs, matching.ID)
			retN++
			if retN >= n {
				break
			}
		}
	}
	return retIDs
}

type AllMatchingLoader struct {
	matchingMap     GroupedMatchings
	cityMatchingMap map[string]CityMatchings

	db       *gorm.DB
	mu       sync.Mutex
	lastLoad time.Time

	cache *cache.Cache

	userMatchingMap map[string]*UserDiscoverMatching
	userMu          sync.Mutex
}

func (l *AllMatchingLoader) GetTopicToMatchings() GroupedMatchings {
	return l.matchingMap
}

func (l *AllMatchingLoader) GetCityToMatchings() map[string]CityMatchings {
	return l.cityMatchingMap
}

var matchingReloadInterval = time.Minute * 5

func (l *AllMatchingLoader) LoadForAnoumynous(ctx context.Context, topicID string, opt UserDiscoverOpt) (retIDs []string) {
	if opt.CityID == "" {
		return l.matchingMap.Load(topicID, opt.Gender.String(), 0, opt.N)
	}
	if l.cityMatchingMap[opt.CityID] == nil {
		return []string{}
	}
	return l.cityMatchingMap[opt.CityID].Load(topicID, opt.Gender.String(), 0, opt.N)
}

func (l *AllMatchingLoader) LoadForUser(ctx context.Context, userID string, topicID string, opt UserDiscoverOpt) (retIDs []string, next string) {
	userDiscoverMatching := l.GenUserDiscoverMatching(ctx, userID)
	matchingLoader := func(topicID, cityID string) []*models.Matching {
		var matchings []*models.Matching
		if cityID == "" {
			matchings = l.matchingMap[topicID]
			if matchings == nil {
				return []*models.Matching{}
			}
		} else {
			cityMatchings := l.cityMatchingMap[cityID]
			if cityMatchings == nil {
				return []*models.Matching{}
			}
			matchings = cityMatchings[topicID]
		}
		ret := make([]*models.Matching, len(matchings))
		copy(ret, matchings)
		lo.Shuffle(ret)
		return ret
	}
	return userDiscoverMatching.LoadMatchingIDs(ctx, matchingLoader, topicID, opt)
}

func (l *AllMatchingLoader) GenUserDiscoverMatching(ctx context.Context, userID string) *UserDiscoverMatching {
	l.userMu.Lock()
	defer l.userMu.Unlock()
	if _, ok := l.userMatchingMap[userID]; !ok {
		l.userMatchingMap[userID] = &UserDiscoverMatching{
			userID:            userID,
			viewedMatchingIDs: map[string]struct{}{},
			matchingMap:       map[string][]*models.Matching{},
		}
	}
	return l.userMatchingMap[userID]
}

func NewAllMatchingLoader(db *gorm.DB) *AllMatchingLoader {
	loader := &AllMatchingLoader{
		db:              db,
		matchingMap:     GroupedMatchings{},
		cityMatchingMap: map[string]CityMatchings{},
		cache:           cache.New(time.Minute, time.Minute),
		userMatchingMap: map[string]*UserDiscoverMatching{},
	}
	UserViewMatching := dbquery.Use(db).UserViewMatching
	MatchingView := dbquery.Use(db).MatchingView

	// 缓存过期后，将用户查看记录写入数据库
	loader.cache.OnEvicted(func(userID string, i interface{}) {
		userDiscoverMatching := i.(*UserDiscoverMatching)
		userDiscoverMatching.mu.Lock()
		defer userDiscoverMatching.mu.Unlock()
		viewedMatchingIDs := lo.Keys(userDiscoverMatching.viewedMatchingIDs)
		if len(viewedMatchingIDs) == 0 {
			return
		}
		// 记录用户查看过的 matching
		err := UserViewMatching.WithContext(context.Background()).Create(&models.UserViewMatching{
			UserID:            userID,
			ViewedMatchingIDs: viewedMatchingIDs,
		})
		if err != nil {
			logger.L.Error("create user view matching failed", zap.Error(err))
		}
		// 更新 matching 被查看的次数
		_, err = MatchingView.WithContext(context.Background()).
			Where(MatchingView.MatchingID.In(viewedMatchingIDs...)).
			UpdateSimple(MatchingView.ViewCount.Add(1))
		if err != nil {
			logger.L.Error("update matching view count failed", zap.Error(err))
		}
	})
	return loader
}

func (l *AllMatchingLoader) Load(ctx context.Context) error {
	if time.Since(l.lastLoad) < matchingReloadInterval {
		return nil
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if time.Since(l.lastLoad) > matchingReloadInterval {
		l.lastLoad = time.Now()
		Matching := dbquery.Use(l.db).Matching
		matchings, err := Matching.WithContext(ctx).Select(
			Matching.ID, Matching.UserID, Matching.CityID, Matching.Gender, Matching.MyGender,
		).Where(Matching.State.Eq(string(models.MatchingStateMatching))).Find()
		if err != nil {
			return nil
		}
		matchingMap := GroupedMatchings{}
		cityMatchingMap := map[string]CityMatchings{}

		for _, m := range matchings {
			matchingMap[m.TopicID] = append(matchingMap[m.TopicID], m)

			cityMatching, ok := cityMatchingMap[m.CityID]
			if !ok {
				cityMatching = CityMatchings{}
			}
			cityMatching[m.TopicID] = append(cityMatching[m.TopicID], m)

			cityMatchingMap[m.CityID] = cityMatching
		}

		l.matchingMap = matchingMap
		l.cityMatchingMap = cityMatchingMap
	}
	return nil
}
