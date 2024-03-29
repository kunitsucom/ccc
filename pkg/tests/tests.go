package tests

import (
	"time"

	"github.com/kunitsucom/ccc/pkg/consts"
	"github.com/kunitsucom/ccc/pkg/domain"
)

// TestDate is
// nolint: gochecknoglobals
var TestDate = time.Date(2022, 2, 22, 22, 22, 22, 222222222, consts.TimeZone("Asia/Tokyo"))

func NewGCPServiceCosts(start time.Time, project, service string, initialCost float64, costChanger float64, currency string, count int) []domain.GCPServiceCost {
	costs := make([]domain.GCPServiceCost, count)

	for i := range costs {
		costs[i] = domain.GCPServiceCost{
			Day:      start.Truncate(24*time.Hour).AddDate(0, 0, i).Format(consts.DateOnly),
			Project:  project,
			Service:  service,
			Cost:     initialCost * costChanger,
			Currency: currency,
		}
	}

	return costs
}
