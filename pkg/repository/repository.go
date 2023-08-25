package repository

import (
	"context"
	"time"

	"github.com/kunitsucom/ccc/pkg/domain"
	"github.com/kunitsucom/ccc/pkg/errors"
	"github.com/kunitsucom/ccc/pkg/repository/bigquery"
	slicez "github.com/kunitsucom/util.go/slices"
)

type Repository struct {
	bigquery *bigquery.BigQuery
}

type Option func(r *Repository) *Repository

func New(opts ...Option) *Repository {
	r := &Repository{}

	for _, opt := range opts {
		r = opt(r)
	}

	return r
}

func WithBigQuery(bigquery *bigquery.BigQuery) Option {
	return func(r *Repository) *Repository {
		r.bigquery = bigquery
		return r
	}
}

func (r *Repository) SUMServiceCostGCPAsc(ctx context.Context, billingTable, billingProject string, from, to time.Time, tz *time.Location, costThreshold float64) ([]domain.GCPServiceCost, error) {
	serviceCostAsc, err := r.bigquery.SUMServiceCostGCPAsc(ctx, billingTable, billingProject, from, to, tz, costThreshold)
	if err != nil {
		return nil, errors.Errorf("(*bigquery.BigQuery).SUMServiceCostGCP: %w", err)
	}

	return serviceCostAsc, nil
}

func (r *Repository) DailyServiceCostGCP(ctx context.Context, billingTable, billingProject string, from, to time.Time, tz *time.Location, costThreshold float64) ([]domain.GCPServiceCost, error) {
	serviceCost, err := r.bigquery.DailyServiceCostGCP(ctx, billingTable, billingProject, from, to, tz, costThreshold)
	if err != nil {
		return nil, errors.Errorf("(*bigquery.BigQuery).DailyServiceCostGCP: %w", err)
	}

	return serviceCost, nil
}

func (r *Repository) DailyServiceCostGCPMapByService(servicesOrderBySUMServiceCostGCP []string, dailyServiceCostGCP []domain.GCPServiceCost) map[string][]domain.GCPServiceCost {
	serviceCost := make(map[string][]domain.GCPServiceCost)
	for _, service := range servicesOrderBySUMServiceCostGCP {
		serviceCost[service] = slicez.Filter(dailyServiceCostGCP, func(index int, source domain.GCPServiceCost) bool {
			// nolint: scopelint
			return service == source.Service
		})
	}

	return serviceCost
}
