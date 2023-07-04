package modelutil_test

import (
	"testing"
	"whale/pkg/models"
	"whale/pkg/modelutil"

	"github.com/stretchr/testify/assert"
)

func TestSimplifyDatePeriod(t *testing.T) {
	ret := modelutil.SimplifyPreferredPeriods([]models.DatePeriod{models.DatePeriodWeekend})
	assert.Len(t, ret, 1)
	assert.Equal(t, models.DatePeriodWeekend.String(), ret[0])

	ret = modelutil.SimplifyPreferredPeriods([]models.DatePeriod{models.DatePeriodWeekendAfternoon, models.DatePeriodWeekendNight})
	assert.Len(t, ret, 1)
	assert.Equal(t, models.DatePeriodWeekend.String(), ret[0])
}
