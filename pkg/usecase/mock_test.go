// nolint: testpackage
package usecase

import (
	"context"
	"io"
	"time"

	"github.com/kunitsucom/ccc/pkg/domain"
)

var _ IRepository = (*repositoryMock)(nil)

// nolint: revive,stylecheck
type repositoryMock struct {
	SUMServiceCostGCPFunc func(ctx context.Context, billingTable string, billingProject string, from time.Time, to time.Time, tz *time.Location, costThreshold float64) ([]domain.GCPServiceCost, error)

	DailyServiceCostGCPFunc func(ctx context.Context, billingTable string, billingProject string, from time.Time, to time.Time, tz *time.Location, costThreshold float64) ([]domain.GCPServiceCost, error)

	DailyServiceCostGCPMapByServiceFunc func(servicesOrderBySUMServiceCostGCP []string, dailyServiceCostGCP []domain.GCPServiceCost) map[string][]domain.GCPServiceCost
}

func (m *repositoryMock) SUMServiceCostGCPAsc(ctx context.Context, billingTable string, billingProject string, from time.Time, to time.Time, tz *time.Location, costThreshold float64) ([]domain.GCPServiceCost, error) {
	return m.SUMServiceCostGCPFunc(ctx, billingTable, billingProject, from, to, tz, costThreshold)
}

func (m *repositoryMock) DailyServiceCostGCP(ctx context.Context, billingTable string, billingProject string, from time.Time, to time.Time, tz *time.Location, costThreshold float64) ([]domain.GCPServiceCost, error) {
	return m.DailyServiceCostGCPFunc(ctx, billingTable, billingProject, from, to, tz, costThreshold)
}

func (m *repositoryMock) DailyServiceCostGCPMapByService(servicesOrderBySUMServiceCostGCP []string, dailyServiceCostGCP []domain.GCPServiceCost) map[string][]domain.GCPServiceCost {
	return m.DailyServiceCostGCPMapByServiceFunc(servicesOrderBySUMServiceCostGCP, dailyServiceCostGCP)
}

var _ IDomain = (*domainMock)(nil)

// nolint: revive,stylecheck
type domainMock struct {
	PlotGraphFunc func(target io.Writer, ps *domain.PlotGraphParameters) error
}

func (m *domainMock) PlotGraph(target io.Writer, ps *domain.PlotGraphParameters) error {
	return m.PlotGraphFunc(target, ps)
}

var _ IInfra = (*infraMock)(nil)

// nolint: revive,stylecheck
type infraMock struct {
	SaveImageFunc func(ctx context.Context, image []byte, imageName string, message string) error
}

func (m *infraMock) SaveImage(ctx context.Context, image []byte, imageName string, message string) error {
	return m.SaveImageFunc(ctx, image, imageName, message)
}
