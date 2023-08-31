package loader

import (
	"context"
	"sync"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/models"

	"github.com/letjoy-club/mida-tool/logger"
	"github.com/letjoy-club/mida-tool/midacontext"
	"github.com/letjoy-club/mida-tool/shortid"
	"github.com/patrickmn/go-cache"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserDiscoverMotion struct {
	userID string

	viewedMotionIDs map[string]struct{}

	mu            sync.Mutex
	motionMap     GroupedMotions
	cityMotionMap map[CityID]CityMotions

	priorityMotion *PriorityMotion
	nextToken      map[string]int
}

type PriorityMotion struct {
	motion     *models.Motion
	categoryID string
}

type UserDiscoverMotionOpt struct {
	CityID   string
	Gender   models.Gender
	TopicIDs []string

	N int

	NextToken string
}

func (u *UserDiscoverMotion) LoadMotionIDs(
	ctx context.Context,
	motionLoader func(topicID, cityID string) []*models.Motion,
	categoryID string,
	opt UserDiscoverMotionOpt,
) (motionIDs []string, next string) {
	u.mu.Lock()
	defer u.mu.Unlock()

	// 如果查看的数量多于 300，就不再返回
	if len(u.viewedMotionIDs) > 300 {
		return []string{}, ""
	}

	startIndex := u.nextToken[opt.NextToken]
	nextIndex := 0

	if opt.CityID == "" {
		_, ok := u.motionMap[categoryID]
		if !ok {
			motions := motionLoader(categoryID, opt.CityID)
			u.motionMap[categoryID] = motions
		}
		motionIDs, nextIndex = u.motionMap.Load(categoryID, opt.TopicIDs, opt.Gender.String(), startIndex, opt.N)
	} else {
		cityMotions, ok := u.cityMotionMap[opt.CityID]
		if !ok {
			u.cityMotionMap[opt.CityID] = CityMotions{}
			cityMotions = u.cityMotionMap[opt.CityID]
		}

		if cityMotions[categoryID] == nil {
			motions := motionLoader(categoryID, opt.CityID)
			u.cityMotionMap[opt.CityID][categoryID] = motions
		}
		motionIDs, nextIndex = cityMotions.Load(categoryID, opt.TopicIDs, opt.Gender.String(), startIndex, opt.N)
	}

	next = shortid.New("", 4)
	u.nextToken[next] = nextIndex

	// 记录已经看过的
	for _, id := range motionIDs {
		u.viewedMotionIDs[id] = struct{}{}
	}
	return motionIDs, next
}

func (u *UserDiscoverMotion) SetPriorityMotion(category string, motion *models.Motion) {
	cityMotion := u.cityMotionMap[motion.CityID]
	if cityMotion == nil {
		return
	}
	categoryMotions := cityMotion[category]
	if categoryMotions == nil {
		return
	}
	if len(categoryMotions) > 0 {
		// 直接覆盖第一个数据
		categoryMotions[0] = motion
	} else {
		categoryMotions = []*models.Motion{motion}
		cityMotion[category] = categoryMotions
	}
}

type CategoryID = string
type GroupedMotions map[CategoryID][]*models.Motion
type CityMotions = GroupedMotions

func (g GroupedMotions) Load(categoryID string, topicIDs []string, gender string, start, n int) ([]string, int) {
	retN := 0
	retIDs := make([]string, 0, n)

	if g[categoryID] == nil {
		return []string{}, 0
	}
	group := g[categoryID]
	cateLen := len(group)
	var topicMap map[string]struct{}
	if len(topicIDs) > 0 {
		topicMap = map[string]struct{}{}
		for _, topicID := range topicIDs {
			topicMap[topicID] = struct{}{}
		}
	}
	for i := start; i < cateLen; i++ {
		motion := group[i]
		if gender != models.GenderN.String() && gender != motion.MyGender {
			continue
		}
		if topicMap != nil {
			if _, ok := topicMap[motion.TopicID]; !ok {
				continue
			}
		}
		retIDs = append(retIDs, motion.ID)
		retN++
		if retN >= n {
			return retIDs, i + 1
		}
	}
	return retIDs, cateLen + 1
}

type CityID = string
type UserID = string

type AllMotionLoader struct {
	categoryMap   GroupedMotions
	cityMotionMap map[CityID]GroupedMotions

	db       *gorm.DB
	mu       sync.Mutex
	lastLoad time.Time

	cache *cache.Cache

	userMu sync.Mutex
}

func (l *AllMotionLoader) GetCategoryToMotions() GroupedMotions {
	return l.categoryMap
}

func (l *AllMotionLoader) GetCityToMotions() map[string]CityMotions {
	return l.cityMotionMap
}

func (l *AllMotionLoader) LoadForAnoumynous(ctx context.Context, categoryID string, opt UserDiscoverMotionOpt) (retIDs []string) {
	if opt.CityID == "" {
		ret, _ := l.categoryMap.Load(categoryID, opt.TopicIDs, opt.Gender.String(), 0, opt.N)
		return ret
	}
	if l.cityMotionMap[opt.CityID] == nil {
		return []string{}
	}
	ret, _ := l.cityMotionMap[opt.CityID].Load(categoryID, opt.TopicIDs, opt.Gender.String(), 0, opt.N)
	return ret
}

func (l *AllMotionLoader) LoadForUser(ctx context.Context, userID string, categoryID string, opt UserDiscoverMotionOpt) (retIDs []string, next string) {
	userDiscoverMotion := l.GenUserDiscoverMotion(ctx, userID)
	matchingLoader := func(categoryID, cityID string) []*models.Motion {
		var motions []*models.Motion
		if cityID == "" {
			motions = l.categoryMap[categoryID]
			if motions == nil {
				return []*models.Motion{}
			}
		} else {
			cityMotions := l.cityMotionMap[cityID]
			if cityMotions == nil {
				return []*models.Motion{}
			}
			motions = cityMotions[categoryID]
		}
		ret := make([]*models.Motion, len(motions))
		copy(ret, motions)
		lo.Shuffle(ret)
		return ret
	}
	return userDiscoverMotion.LoadMotionIDs(ctx, matchingLoader, categoryID, opt)
}

func (l *AllMotionLoader) GenUserDiscoverMotion(ctx context.Context, userID string) *UserDiscoverMotion {
	l.userMu.Lock()
	defer l.userMu.Unlock()
	userDiscoverMotion, ok := l.cache.Get(userID)
	if !ok {
		userDiscoverMotion = &UserDiscoverMotion{
			userID:          userID,
			viewedMotionIDs: map[string]struct{}{},
			cityMotionMap:   map[CityID]CityMotions{},
			motionMap:       map[CategoryID][]*models.Motion{},
			nextToken:       map[string]int{},
		}
	}
	l.cache.SetDefault(userID, userDiscoverMotion)
	return userDiscoverMotion.(*UserDiscoverMotion)
}

func NewAllMotionLoader(db *gorm.DB) *AllMotionLoader {
	loader := &AllMotionLoader{
		db:            db,
		categoryMap:   map[CategoryID][]*models.Motion{},
		cityMotionMap: map[CityID]CityMotions{},
		cache:         cache.New(time.Minute, time.Minute),
	}
	MotionViewHistory := dbquery.Use(db).MotionViewHistory
	Motion := dbquery.Use(db).Motion

	// 缓存过期后，将用户查看记录写入数据库
	loader.cache.OnEvicted(func(userID string, i interface{}) {
		userDiscoverMotion := i.(*UserDiscoverMotion)
		userDiscoverMotion.mu.Lock()
		defer userDiscoverMotion.mu.Unlock()
		viewedMotionIDs := lo.Keys(userDiscoverMotion.viewedMotionIDs)
		if len(viewedMotionIDs) == 0 {
			return
		}
		// 记录用户查看过的 matching
		err := MotionViewHistory.WithContext(context.Background()).Create(&models.MotionViewHistory{
			UserID:    userID,
			MotionIDs: viewedMotionIDs,
		})
		if err != nil {
			logger.L.Error("create user view motions failed", zap.Error(err))
		}
		// 更新 matching 被查看的次数
		_, err = Motion.WithContext(context.Background()).
			Where(Motion.ID.In(viewedMotionIDs...)).
			UpdateSimple(Motion.ViewCount.Add(1))
		if err != nil {
			logger.L.Error("update matching view count failed", zap.Error(err))
		}
	})
	return loader
}

func (l *AllMotionLoader) SetLatestData(ctx context.Context, userID string, motion *models.Motion) {
	discoverMotionInterface, ok := l.cache.Get(userID)
	if !ok {
		return
	}
	discoverMotion := discoverMotionInterface.(*UserDiscoverMotion)
	topicCategory := midacontext.GetLoader[Loader](ctx).TopicCategory
	topicCategory.Load(ctx)
	categoryID, ok := topicCategory.topic2category[motion.TopicID]
	if !ok {
		return
	}
	discoverMotion.SetPriorityMotion(categoryID, motion)
}

func (l *AllMotionLoader) Load(ctx context.Context) error {
	if time.Since(l.lastLoad) < matchingReloadInterval {
		return nil
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if time.Since(l.lastLoad) > matchingReloadInterval {
		l.lastLoad = time.Now()
		Motion := dbquery.Use(l.db).Motion
		motions, err := Motion.WithContext(ctx).Select(
			Motion.ID, Motion.UserID, Motion.CityID, Motion.Gender, Motion.MyGender, Motion.TopicID,
		).Where(Motion.Active.Is(true)).Find()
		if err != nil {
			logger.L.Error("load motions failed", zap.Error(err))
			return nil
		}
		categoryMap := GroupedMotions{}
		citiesMotionMap := map[CityID]CityMotions{}

		topicCategory := midacontext.GetLoader[Loader](ctx).TopicCategory
		topicCategory.Load(ctx)

		for _, m := range motions {
			categoryID := topicCategory.Category(m.TopicID)
			categoryMap[categoryID] = append(categoryMap[categoryID], m)

			cityMotion, ok := citiesMotionMap[m.CityID]
			if !ok {
				cityMotion = CityMotions{}
			}
			cityMotion[categoryID] = append(cityMotion[categoryID], m)
			citiesMotionMap[m.CityID] = cityMotion
		}

		l.categoryMap = categoryMap
		l.cityMotionMap = citiesMotionMap
	}
	return nil
}

func insertIfNotExist() {
}
