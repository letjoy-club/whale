package loader

import (
	"context"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/models"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/letjoy-club/mida-tool/loaderutil"
	"gorm.io/gorm"
)

func NewMatchingOfferSummaryLoader(db *gorm.DB) *dataloader.Loader[string, *models.MatchingOfferSummary] {
	MatchingOffer := dbquery.Use(db).MatchingOfferSummary
	return loaderutil.NewItemLoader(db, func(ctx context.Context, ids []string) ([]*models.MatchingOfferSummary, error) {
		matchingOffers, err := MatchingOffer.WithContext(ctx).Where(MatchingOffer.MatchingID.In(ids...)).Find()
		if err != nil {
			return nil, err
		}
		return matchingOffers, nil
	}, func(m map[string]*models.MatchingOfferSummary, v *models.MatchingOfferSummary) {
		m[v.MatchingID] = v
	}, time.Minute*5, loaderutil.Placeholder(func(ctx context.Context, id string) (*models.MatchingOfferSummary, error) {
		return &models.MatchingOfferSummary{MatchingID: id}, nil
	}),
	)
}

type MatchingOffers struct {
	matchingID string

	records []*models.MatchingOfferRecord

	unprocessCount *int
}

func (m *MatchingOffers) OfferRecords() []*models.MatchingOfferRecord {
	return m.records
}

func (m *MatchingOffers) UnprocessCount() int {
	if m.unprocessCount != nil {
		return *m.unprocessCount
	}
	var count int
	for _, record := range m.records {
		if record.State == string(models.MatchingOfferStateUnprocessed) {
			count++
		}
	}
	m.unprocessCount = &count
	return count

}

func NewOutMatchingOfferLoader(db *gorm.DB) *dataloader.Loader[string, *MatchingOffers] {
	MatchingOffer := dbquery.Use(db).MatchingOfferRecord
	return loaderutil.NewAggregatorLoader(db, func(ctx context.Context, ids []string) ([]*models.MatchingOfferRecord, error) {
		matchingOffers, err := MatchingOffer.WithContext(ctx).Where(MatchingOffer.MatchingID.In(ids...)).Find()
		if err != nil {
			return nil, err
		}
		return matchingOffers, nil
	}, func(m map[string]*MatchingOffers, v *models.MatchingOfferRecord) {
		if _, ok := m[v.MatchingID]; !ok {
			m[v.MatchingID] = &MatchingOffers{matchingID: v.MatchingID, records: []*models.MatchingOfferRecord{v}}
		} else {
			m[v.MatchingID].records = append(m[v.MatchingID].records, v)
		}
	}, time.Minute*5, loaderutil.Placeholder(func(ctx context.Context, id string) (*MatchingOffers, error) {
		return &MatchingOffers{matchingID: id, records: []*models.MatchingOfferRecord{}}, nil
	}),
	)
}

func NewInMatchingOfferLoader(db *gorm.DB) *dataloader.Loader[string, *MatchingOffers] {
	MatchingOffer := dbquery.Use(db).MatchingOfferRecord
	return loaderutil.NewAggregatorLoader(db, func(ctx context.Context, ids []string) ([]*models.MatchingOfferRecord, error) {
		matchingOffers, err := MatchingOffer.WithContext(ctx).Where(MatchingOffer.ToMatchingID.In(ids...)).Find()
		if err != nil {
			return nil, err
		}
		return matchingOffers, nil
	}, func(m map[string]*MatchingOffers, v *models.MatchingOfferRecord) {
		if _, ok := m[v.ToMatchingID]; !ok {
			m[v.ToMatchingID] = &MatchingOffers{matchingID: v.ToMatchingID, records: []*models.MatchingOfferRecord{v}}
		} else {
			m[v.ToMatchingID].records = append(m[v.ToMatchingID].records, v)
		}
	}, time.Minute*5, loaderutil.Placeholder(func(ctx context.Context, id string) (*MatchingOffers, error) {
		return &MatchingOffers{matchingID: id, records: []*models.MatchingOfferRecord{}}, nil
	}),
	)
}
