package loader

import (
	"context"
	"math/rand"
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
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

type UserDiscoverMotion struct {
	userID string

	viewedMotionIDs map[string]struct{}

	mu            sync.Mutex
	motionMap     GroupedMotions
	cityMotionMap map[CityID]CityMotions

	nextToken map[string]int
}

type UserDiscoverMotionOpt struct {
	CityID   string
	Gender   models.Gender
	TopicIDs []string

	N int

	NextToken string
	LastID    string

	CategoryID string

	Type models.MotionType
}

var sessionMaxViewed = 300

func (u *UserDiscoverMotion) Viewed(motionIDs []string) int {
	u.mu.Lock()
	defer u.mu.Unlock()
	// 如果查看的数量多于 300，就不再返回
	if len(u.viewedMotionIDs) > sessionMaxViewed {
		return 0
	}
	for _, id := range motionIDs {
		u.viewedMotionIDs[id] = struct{}{}
	}
	return len(u.viewedMotionIDs)
}

func (u *UserDiscoverMotion) LoadMotionIDs(
	ctx context.Context,
	motionLoader func(topicID, cityID string) []*models.Motion,
	opt UserDiscoverMotionOpt,
) (motionIDs []string, next string) {
	u.mu.Lock()
	defer u.mu.Unlock()

	// 如果查看的数量多于 300，就不再返回
	if len(u.viewedMotionIDs) > sessionMaxViewed {
		return []string{}, ""
	}

	startIndex := u.nextToken[opt.NextToken]
	nextIndex := 0

	if opt.CityID == "" {
		_, ok := u.motionMap[opt.CategoryID]
		if !ok {
			motions := motionLoader(opt.CategoryID, opt.CityID)
			u.motionMap[opt.CategoryID] = motions
		}
		motionIDs, nextIndex = u.motionMap.Load(opt.CategoryID, opt.TopicIDs, opt.Type, opt.Gender.String(), startIndex, opt.N)
	} else {
		cityMotions, ok := u.cityMotionMap[opt.CityID]
		if !ok {
			u.cityMotionMap[opt.CityID] = CityMotions{}
			cityMotions = u.cityMotionMap[opt.CityID]
		}

		if cityMotions[opt.CategoryID] == nil {
			motions := motionLoader(opt.CategoryID, opt.CityID)
			u.cityMotionMap[opt.CityID][opt.CategoryID] = motions
		}
		motionIDs, nextIndex = cityMotions.Load(opt.CategoryID, opt.TopicIDs, opt.Type, opt.Gender.String(), startIndex, opt.N)
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

func (g GroupedMotions) Load(categoryID string, topicIDs []string, motionType models.MotionType, gender string, start, n int) ([]string, int) {
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

		if motionType != models.MotionTypeGirlOnly {
			// 不是女生专区的，需要遵从性别设置
			if gender != models.GenderN.String() && gender != motion.MyGender {
				continue
			}
		}
		if topicMap != nil {
			if _, ok := topicMap[motion.TopicID]; !ok {
				continue
			}
		}

		// 如果是请求特定类型 motion
		if motionType != models.MotionTypeAll {
			t := GetMotionType(motion)
			if motionType != t {
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

	// 当有新的 motion 时，需要更新这个字段
	latestMotions []*models.Motion
}

// GetOrderedMotions 获取按照时间排序的 motion id
func (l *AllMotionLoader) GetOrderedMotions(ctx context.Context, opt UserDiscoverMotionOpt) []string {
	latestMotions := l.latestMotions

	index := len(latestMotions)
	var exist bool

	if opt.LastID != "" {
		index, exist = slices.BinarySearchFunc(latestMotions, opt.LastID, func(m *models.Motion, olderThanID string) int {
			if m.ID < olderThanID {
				return -1
			}
			if m.ID == olderThanID {
				return 0
			}
			return 1
		})

		if !exist {
			return []string{}
		}
	}

	topicCategory := midacontext.GetLoader[Loader](ctx).TopicCategory
	topicCategory.Load(ctx)

	topicIDs := opt.TopicIDs
	var topicMap map[string]struct{}
	if len(topicIDs) > 0 {
		topicMap = map[string]struct{}{}
		for _, topicID := range topicIDs {
			topicMap[topicID] = struct{}{}
		}
	}

	ret := make([]string, 0, opt.N)
	retN := opt.N
	for i := index - 1; i >= 0; i-- {
		motion := latestMotions[i]
		if opt.CityID != motion.CityID {
			continue
		}
		if opt.CategoryID != AllCategoryID {
			category := topicCategory.Category(motion.TopicID)
			if opt.CategoryID != category {
				continue
			}
		}
		if topicIDs != nil {
			// 指定了 topicIDs，进行过滤
			if _, ok := topicMap[motion.TopicID]; !ok {
				continue
			}
		}
		// 如果是请求女生专区, 忽略性别设置
		if opt.Type != models.MotionTypeGirlOnly {
			if opt.Gender != models.GenderN && opt.Gender.String() != motion.MyGender {
				// 根据用户期望性别进行过滤
				continue
			}
		}
		// 如果是请求特定类型 motion
		if opt.Type != models.MotionTypeAll {
			t := GetMotionType(motion)
			if opt.Type != t {
				continue
			}
		}
		ret = append(ret, motion.ID)
		retN--
		if retN <= 0 {
			break
		}
	}
	return ret
}

func (l *AllMotionLoader) GetCategoryToMotions() GroupedMotions {
	return l.categoryMap
}

func (l *AllMotionLoader) GetCityToMotions() map[string]CityMotions {
	return l.cityMotionMap
}

func (l *AllMotionLoader) LoadForAnoumynous(ctx context.Context, opt UserDiscoverMotionOpt) (retIDs []string) {
	if opt.CityID == "" {
		ret, _ := l.categoryMap.Load(opt.CategoryID, opt.TopicIDs, opt.Type, opt.Gender.String(), 0, opt.N)
		return ret
	}
	if l.cityMotionMap[opt.CityID] == nil {
		return []string{}
	}
	ret, _ := l.cityMotionMap[opt.CityID].Load(opt.CategoryID, opt.TopicIDs, opt.Type, opt.Gender.String(), 0, opt.N)
	return ret
}

func (l *AllMotionLoader) LoadForUser(ctx context.Context, userID string, opt UserDiscoverMotionOpt) (retIDs []string, next string) {
	userDiscoverMotion := l.GenUserDiscoverMotion(ctx, userID)
	motionLoader := func(categoryID, cityID string) []*models.Motion {
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
		length := len(motions)
		ret := make([]*models.Motion, length)
		copy(ret, motions)
		lo.Shuffle(ret)

		// 按照一定优先级排序，排序两次，让最近创建的，点赞数高的尽量靠前（排序次数越多，这个趋势越明显）
		PrioritySort(ret)
		PrioritySort(ret)
		return ret
	}
	return userDiscoverMotion.LoadMotionIDs(ctx, motionLoader, opt)
}

// 按照创建时间，点赞量，以一定概率排序
func PrioritySort(motions []*models.Motion) {
	length := len(motions)
	for i := 0; i < length; i++ {
		curr := motions[i]
		targetIndex := i + rand.Intn(length-i)
		target := motions[targetIndex]
		if curr.CreatedAt.Before(target.CreatedAt) {
			// 如果当前的创建时间比后面的某一个 motion 的创建时间早，那么交换（让最近创建的尽量靠前）
			motions[i], motions[targetIndex] = motions[targetIndex], motions[i]
		} else if curr.ThumbsUpCount < target.ThumbsUpCount {
			// 相比时间，点赞数据重要性排第二
			// 如果当前的点赞量比后面的某一个 motion 的点赞量少，那么交换（让点赞量多的尽量靠前）
			motions[i], motions[targetIndex] = motions[targetIndex], motions[i]
		}
	}
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
		motions, err := Motion.WithContext(ctx).
			Select(
				Motion.ID, Motion.UserID,
				// 用于筛选
				Motion.CityID, Motion.Gender, Motion.MyGender, Motion.TopicID, Motion.Quick,
				// 点赞数，时间用于排序
				Motion.ThumbsUpCount, Motion.CreatedAt,
			).
			Where(Motion.Active.Is(true)).
			// 默认按创建时间倒序
			Order(Motion.ID).Find()
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

		categoryMap[AllCategoryID] = motions

		l.categoryMap = categoryMap
		l.cityMotionMap = citiesMotionMap
		l.latestMotions = motions
	}
	return nil
}

var AllCategoryID = "All"

// func creativeSort(motions []*models.Motion) {
// 	length := len(motions)
// 	for i := 0; i < length; i++ {
// 		curr := motions[i]
// 		targetIndex := rand.Intn(length-i) + i
// 		target := motions[i]
// 		if curr.ThumbsUpCount < target.ThumbsUpCount || curr.CreatedAt.Before(target.CreatedAt) {
// 			motions[i], motions[targetIndex] = motions[targetIndex], motions[i]
// 		}
// 	}
// }

func (l *AllMotionLoader) AppendNewMotion(ctx context.Context, motion *models.Motion) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.latestMotions = append(l.latestMotions, motion)
}

func GetMotionType(m *models.Motion) models.MotionType {
	if m.Quick {
		return models.MotionTypeQuick
	}
	if m.MyGender == models.GenderF.String() && m.Gender == m.MyGender {
		return models.MotionTypeGirlOnly
	}
	return models.MotionTypeNormal
}
