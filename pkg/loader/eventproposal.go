package loader

import (
	"sync"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/models"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/letjoy-club/mida-tool/dbutil"
	"github.com/letjoy-club/mida-tool/loaderutil"
	"github.com/letjoy-club/mida-tool/midacontext"
	"golang.org/x/exp/slices"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

func NewEventProposalLoader(db *gorm.DB) *dataloader.Loader[string, *models.EventProposal] {
	EventProposal := dbquery.Use(db).EventProposal
	return loaderutil.NewItemLoader(db, func(ctx context.Context, keys []string) ([]*models.EventProposal, error) {
		events, err := EventProposal.WithContext(ctx).Where(EventProposal.ID.In(keys...)).Find()
		return events, err
	}, func(m map[string]*models.EventProposal, v *models.EventProposal) {
		m[v.ID] = v
	}, time.Minute)
}

type EventProposalParticipants struct {
	ID string

	Events []*models.UserJoinEventProposal

	activeUserIDs []string
}

func (e *EventProposalParticipants) ActiveEvents() []*models.UserJoinEventProposal {
	events := make([]*models.UserJoinEventProposal, 0, len(e.Events))
	for _, e := range e.Events {
		if e.State == models.JoinEventStateJoined.String() {
			events = append(events, e)
		}
	}
	return events
}

func (e *EventProposalParticipants) ActiveUserIDs() []string {
	if e.activeUserIDs == nil {
		ids := make([]string, 0, len(e.activeUserIDs))
		for _, e := range e.Events {
			if e.State == models.JoinEventStateJoined.String() {
				ids = append(ids, e.UserID)
			}
		}
		e.activeUserIDs = ids
	}
	return e.activeUserIDs
}

func (e EventProposalParticipants) HasActiveUser(uid string) bool {
	for _, e := range e.Events {
		if e.UserID == uid && e.State == string(models.JoinEventStateJoined) {
			return true
		}
	}
	return false
}

func (e EventProposalParticipants) HasUser(uid string) bool {
	for _, e := range e.Events {
		if e.UserID == uid {
			return true
		}
	}
	return false
}

func (e EventProposalParticipants) GetEvent(uid string) *models.UserJoinEventProposal {
	for _, e := range e.Events {
		if e.UserID == uid {
			return e
		}
	}
	return nil
}

func NewEventProposalParticipantsLoader(db *gorm.DB) *dataloader.Loader[string, *EventProposalParticipants] {
	UserJoinEventProposal := dbquery.Use(db).UserJoinEventProposal
	return loaderutil.NewAggregatorLoader(db, func(ctx context.Context, keys []string) ([]*models.UserJoinEventProposal, error) {
		events, err := UserJoinEventProposal.WithContext(ctx).Where(UserJoinEventProposal.EventID.In(keys...)).Find()
		return events, err
	}, func(m map[string]*EventProposalParticipants, v *models.UserJoinEventProposal) {
		if participant, ok := m[v.EventID]; ok {
			participant.Events = append(participant.Events, v)
		} else {
			m[v.EventID] = &EventProposalParticipants{
				ID:     v.EventID,
				Events: []*models.UserJoinEventProposal{v},
			}
		}
	}, time.Minute, loaderutil.Placeholder(func(ctx context.Context, id string) (*EventProposalParticipants, error) {
		return &EventProposalParticipants{
			ID:     id,
			Events: []*models.UserJoinEventProposal{},
		}, nil
	}))
}

type AllEventProposalLoader struct {
	lastUpdate time.Time

	mu sync.Mutex

	events []*models.EventProposal
}

func NewAllEventProposalLoader() *AllEventProposalLoader {
	return &AllEventProposalLoader{
		events: []*models.EventProposal{},
	}
}

func (a *AllEventProposalLoader) AppendNewData(e *models.EventProposal) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.events = append(a.events, e)
}

func (a *AllEventProposalLoader) Load(ctx context.Context) error {
	if a.lastUpdate.After(time.Now().Add(time.Minute * 5)) {
		return nil
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	// double check
	if a.lastUpdate.After(time.Now().Add(time.Minute * 5)) {
		return nil
	}

	db := dbutil.GetDB(ctx)
	EventProposal := dbquery.Use(db).EventProposal
	events, err := EventProposal.WithContext(ctx).Where(EventProposal.Active.Is(true)).Find()
	if err != nil {
		return err
	}

	a.events = events
	a.lastUpdate = time.Now()
	return nil
}

func (a *AllEventProposalLoader) GetOrdered(ctx context.Context, categoryID string, topicIDs []string, cityID, lastID string, n int) []string {
	events := a.events

	var index int
	var exist bool

	if lastID != "" {
		index, exist = slices.BinarySearchFunc(events, lastID, func(e *models.EventProposal, olderThanID string) int {
			if e.ID < olderThanID {
				return -1
			}
			if e.ID == olderThanID {
				return 0
			}
			return 1
		})

		if !exist {
			return []string{}
		}
	} else {
		index = len(events)
	}

	category := midacontext.GetLoader[Loader](ctx).TopicCategory
	category.Load(ctx)

	var topicMap map[string]struct{}
	if len(topicIDs) > 0 {
		topicMap = make(map[string]struct{}, len(topicIDs))
		for _, id := range topicIDs {
			topicMap[id] = struct{}{}
		}
	}

	ids := make([]string, 0, n)
	for i := index - 1; i >= 0; i-- {
		event := events[i]

		if event.CityID != cityID {
			continue
		}

		if categoryID != AllCategoryID {
			if categoryID != category.Category(event.TopicID) {
				continue
			}
		}

		if topicMap != nil {
			if _, exist := topicMap[event.TopicID]; !exist {
				continue
			}
		}

		ids = append(ids, event.ID)
	}
	return ids
}
