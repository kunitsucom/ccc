package usecase

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/kunitsuinc/ccc/pkg/domain"
	"github.com/kunitsuinc/ccc/pkg/infra"
	"github.com/kunitsuinc/ccc/pkg/repository"
)

var ErrMixedCurrenciesDataSourceIsNotSupported = errors.New("usecase: mixed currencies data source is not supported")

type UseCase struct {
	repository RepositoryIF
	domain     DomainIF
	infra      InfraIF
}

type Option func(r *UseCase) *UseCase

func New(opts ...Option) *UseCase {
	u := &UseCase{}

	for _, opt := range opts {
		u = opt(u)
	}

	return u
}

var _ RepositoryIF = (*repository.Repository)(nil)

type RepositoryIF interface {
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

var _ DomainIF = (*domain.Domain)(nil)

type DomainIF interface {
	PlotGraph(target io.Writer, ps *domain.PlotGraphParameters) error
}

func WithDomain(d *domain.Domain) Option {
	return func(u *UseCase) *UseCase {
		u.domain = d
		return u
	}
}

type InfraIF interface {
	SaveImage(ctx context.Context, image io.ReadSeeker, imageName string, message string) error
}

func WithInfra(i *infra.Infra) Option {
	return func(u *UseCase) *UseCase {
		u.infra = i
		return u
	}
}
