package usecase

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/kunitsucom/ccc/pkg/domain"
	"github.com/kunitsucom/ccc/pkg/infra"
	"github.com/kunitsucom/ccc/pkg/repository"
)

var ErrMixedCurrenciesDataSourceIsNotSupported = errors.New("usecase: mixed currencies data source is not supported")

type UseCase struct {
	repository IRepository
	domain     IDomain
	infra      IInfra
}

type Option func(r *UseCase) *UseCase

func New(opts ...Option) *UseCase {
	u := &UseCase{}

	for _, opt := range opts {
		u = opt(u)
	}

	return u
}

var _ IRepository = (*repository.Repository)(nil)

type IRepository interface {
	SUMServiceCostGCPAsc(ctx context.Context, billingTable string, billingProject string, from time.Time, to time.Time, tz *time.Location, costThreshold float64) ([]domain.GCPServiceCost, error)
	DailyServiceCostGCP(ctx context.Context, billingTable string, billingProject string, from time.Time, to time.Time, tz *time.Location, costThreshold float64) ([]domain.GCPServiceCost, error)
	DailyServiceCostGCPMapByService(servicesOrderBySUMServiceCostGCP []string, dailyServiceCostGCP []domain.GCPServiceCost) map[string][]domain.GCPServiceCost
}

func WithRepository(r *repository.Repository) Option {
	return func(u *UseCase) *UseCase {
		u.repository = r
		return u
	}
}

var _ IDomain = (*domain.Domain)(nil)

type IDomain interface {
	PlotGraph(target io.Writer, ps *domain.PlotGraphParameters) error
}

func WithDomain(d *domain.Domain) Option {
	return func(u *UseCase) *UseCase {
		u.domain = d
		return u
	}
}

type IInfra interface {
	SaveImage(ctx context.Context, image []byte, imageName string, message string) error
}

func WithInfra(i *infra.Infra) Option {
	return func(u *UseCase) *UseCase {
		u.infra = i
		return u
	}
}
