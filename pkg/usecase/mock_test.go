// nolint: testpackage
package usecase

import (
	"context"
	"io"
	"time"

	"github.com/kunitsuinc/ccc/pkg/domain"
)

var _ RepositoryIF = (*repositoryMock)(nil)

// nolint: revive,stylecheck
type repositoryMock struct {
	SUMServiceCostGCP_GCPServiceCost []domain.GCPServiceCost
	SUMServiceCostGCP_error          error

	DailyServiceCostGCP_GCPServiceCost []domain.GCPServiceCost
	DailyServiceCostGCP_error          error

	DailyServiceCostGCPMapByService_map_string_GCPServiceCost map[string][]domain.GCPServiceCost
}

func (m *repositoryMock) SUMServiceCostGCPAsc(ctx context.Context, billingTable string, billingProject string, from time.Time, to time.Time, tz *time.Location, costThreshold float64) ([]domain.GCPServiceCost, error) {
	return m.SUMServiceCostGCP_GCPServiceCost, m.SUMServiceCostGCP_error
}

func (m *repositoryMock) DailyServiceCostGCP(ctx context.Context, billingTable string, billingProject string, from time.Time, to time.Time, tz *time.Location, costThreshold float64) ([]domain.GCPServiceCost, error) {
	return m.DailyServiceCostGCP_GCPServiceCost, m.DailyServiceCostGCP_error
}

func (m *repositoryMock) DailyServiceCostGCPMapByService(servicesOrderBySUMServiceCostGCP []string, dailyServiceCostGCP []domain.GCPServiceCost) map[string][]domain.GCPServiceCost {
	return m.DailyServiceCostGCPMapByService_map_string_GCPServiceCost
}

var _ DomainIF = (*domainMock)(nil)

// nolint: revive,stylecheck
type domainMock struct {
	PlotGraph_error error
}

func (m *domainMock) PlotGraph(target io.Writer, ps *domain.PlotGraphParameters) error {
	return m.PlotGraph_error
}

var _ InfraIF = (*infraMock)(nil)

// nolint: revive,stylecheck
type infraMock struct {
	SaveImage_error error
}

func (m *infraMock) SaveImage(ctx context.Context, image io.ReadSeeker, imageName string, message string) error {
	return m.SaveImage_error
}
