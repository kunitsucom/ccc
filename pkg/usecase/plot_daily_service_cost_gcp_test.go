// nolint: testpackage
package usecase

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/kunitsucom/ccc/pkg/domain"
	"github.com/kunitsucom/ccc/pkg/errors"
	"github.com/kunitsucom/ccc/pkg/tests"
	errorz "github.com/kunitsucom/util.go/errors"
	testz "github.com/kunitsucom/util.go/test"
)

func TestUsecase_PlotDailyServiceCostGCP(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		u := &UseCase{
			repository: &repositoryMock{
				SUMServiceCostGCPFunc: func(ctx context.Context, billingTable string, billingProject string, from time.Time, to time.Time, tz *time.Location, costThreshold float64) ([]domain.GCPServiceCost, error) {
					return tests.NewGCPServiceCosts(tests.TestDate, "test-project", "TestService", 123.45, 1, "USD", 5), nil
				},
				DailyServiceCostGCPFunc: func(ctx context.Context, billingTable string, billingProject string, from time.Time, to time.Time, tz *time.Location, costThreshold float64) ([]domain.GCPServiceCost, error) {
					return tests.NewGCPServiceCosts(tests.TestDate, "test-project", "TestService", 123.45, 1, "USD", 5), nil
				},
				DailyServiceCostGCPMapByServiceFunc: func(servicesOrderBySUMServiceCostGCP []string, dailyServiceCostGCP []domain.GCPServiceCost) map[string][]domain.GCPServiceCost {
					return map[string][]domain.GCPServiceCost{"TestService": tests.NewGCPServiceCosts(tests.TestDate, "test-project", "TestService", 123.45, 1, "USD", 5)}
				},
			},
			domain: &domainMock{
				PlotGraphFunc: func(target io.Writer, ps *domain.PlotGraphParameters) error { return nil },
			},
			infra: &infraMock{
				SaveImageFunc: func(ctx context.Context, image []byte, imageName string, message string) error { return nil },
			},
		}
		ctx := context.Background()
		buf := bytes.NewBuffer(nil)
		err := u.PlotDailyServiceCostGCP(ctx, buf, &PlotDailyServiceCostGCPParameters{})
		if err != nil {
			t.Errorf("err != nil: %v", err)
		}
	})

	t.Run("failure(SUMServiceCostGCP)", func(t *testing.T) {
		t.Parallel()
		u := &UseCase{
			repository: &repositoryMock{
				SUMServiceCostGCPFunc: func(ctx context.Context, billingTable string, billingProject string, from time.Time, to time.Time, tz *time.Location, costThreshold float64) ([]domain.GCPServiceCost, error) {
					return nil, testz.ErrTestError
				},
				DailyServiceCostGCPFunc: func(ctx context.Context, billingTable string, billingProject string, from time.Time, to time.Time, tz *time.Location, costThreshold float64) ([]domain.GCPServiceCost, error) {
					return tests.NewGCPServiceCosts(tests.TestDate, "test-project", "TestService", 123.45, 1, "USD", 5), nil
				},
				DailyServiceCostGCPMapByServiceFunc: func(servicesOrderBySUMServiceCostGCP []string, dailyServiceCostGCP []domain.GCPServiceCost) map[string][]domain.GCPServiceCost {
					return map[string][]domain.GCPServiceCost{"TestService": tests.NewGCPServiceCosts(tests.TestDate, "test-project", "TestService", 123.45, 1, "USD", 5)}
				},
			},
		}
		ctx := context.Background()
		buf := bytes.NewBuffer(nil)
		err := u.PlotDailyServiceCostGCP(ctx, buf, &PlotDailyServiceCostGCPParameters{})
		if !errorz.Contains(err, "(IRepository).SUMServiceCostGCP") {
			t.Errorf("err not contain (IRepository).SUMServiceCostGCP: %v", err)
		}
	})

	t.Run("failure(DailyServiceCostGCP)", func(t *testing.T) {
		t.Parallel()
		u := &UseCase{
			repository: &repositoryMock{
				SUMServiceCostGCPFunc: func(ctx context.Context, billingTable string, billingProject string, from time.Time, to time.Time, tz *time.Location, costThreshold float64) ([]domain.GCPServiceCost, error) {
					return tests.NewGCPServiceCosts(tests.TestDate, "test-project", "TestService", 123.45, 1, "USD", 5), nil
				},
				DailyServiceCostGCPFunc: func(ctx context.Context, billingTable string, billingProject string, from time.Time, to time.Time, tz *time.Location, costThreshold float64) ([]domain.GCPServiceCost, error) {
					return nil, testz.ErrTestError
				},
			},
		}
		ctx := context.Background()
		buf := bytes.NewBuffer(nil)
		err := u.PlotDailyServiceCostGCP(ctx, buf, &PlotDailyServiceCostGCPParameters{})
		if !errorz.Contains(err, "(IRepository).DailyServiceCostGCP") {
			t.Errorf("err not contain (IRepository).DailyServiceCostGCP: %v", err)
		}
	})

	t.Run("failure(ErrMixedCurrenciesDataSourceIsNotSupported)", func(t *testing.T) {
		t.Parallel()
		u := &UseCase{
			repository: &repositoryMock{
				SUMServiceCostGCPFunc: func(ctx context.Context, billingTable string, billingProject string, from time.Time, to time.Time, tz *time.Location, costThreshold float64) ([]domain.GCPServiceCost, error) {
					return append(tests.NewGCPServiceCosts(tests.TestDate, "test-project", "TestService", 123.45, 1, "USD", 5), tests.NewGCPServiceCosts(tests.TestDate, "test-project", "TestService", 123.45, 1, "JPY", 5)...), nil
				},
				DailyServiceCostGCPFunc: func(ctx context.Context, billingTable string, billingProject string, from time.Time, to time.Time, tz *time.Location, costThreshold float64) ([]domain.GCPServiceCost, error) {
					return append(tests.NewGCPServiceCosts(tests.TestDate, "test-project", "TestService", 123.45, 1, "USD", 5), tests.NewGCPServiceCosts(tests.TestDate, "test-project", "TestService", 123.45, 1, "JPY", 5)...), nil
				},
				DailyServiceCostGCPMapByServiceFunc: func(servicesOrderBySUMServiceCostGCP []string, dailyServiceCostGCP []domain.GCPServiceCost) map[string][]domain.GCPServiceCost {
					return map[string][]domain.GCPServiceCost{"TestService": tests.NewGCPServiceCosts(tests.TestDate, "test-project", "TestService", 123.45, 1, "USD", 5)}
				},
			},
		}
		ctx := context.Background()
		buf := bytes.NewBuffer(nil)
		err := u.PlotDailyServiceCostGCP(ctx, buf, &PlotDailyServiceCostGCPParameters{})
		if !errors.Is(err, ErrMixedCurrenciesDataSourceIsNotSupported) {
			t.Errorf("err != ErrMixedCurrenciesDataSourceIsNotSupported: %v", err)
		}
	})

	t.Run("failure(PlotGraph)", func(t *testing.T) {
		t.Parallel()
		u := &UseCase{
			repository: &repositoryMock{
				SUMServiceCostGCPFunc: func(ctx context.Context, billingTable string, billingProject string, from time.Time, to time.Time, tz *time.Location, costThreshold float64) ([]domain.GCPServiceCost, error) {
					return tests.NewGCPServiceCosts(tests.TestDate, "test-project", "TestService", 123.45, 1, "USD", 5), nil
				},
				DailyServiceCostGCPFunc: func(ctx context.Context, billingTable string, billingProject string, from time.Time, to time.Time, tz *time.Location, costThreshold float64) ([]domain.GCPServiceCost, error) {
					return tests.NewGCPServiceCosts(tests.TestDate, "test-project", "TestService", 123.45, 1, "USD", 5), nil
				},
				DailyServiceCostGCPMapByServiceFunc: func(servicesOrderBySUMServiceCostGCP []string, dailyServiceCostGCP []domain.GCPServiceCost) map[string][]domain.GCPServiceCost {
					return map[string][]domain.GCPServiceCost{"TestService": tests.NewGCPServiceCosts(tests.TestDate, "test-project", "TestService", 123.45, 1, "USD", 5)}
				},
			},
			domain: &domainMock{
				PlotGraphFunc: func(target io.Writer, ps *domain.PlotGraphParameters) error { return testz.ErrTestError },
			},
		}
		ctx := context.Background()
		buf := bytes.NewBuffer(nil)
		err := u.PlotDailyServiceCostGCP(ctx, buf, &PlotDailyServiceCostGCPParameters{})
		if !errorz.Contains(err, "(IDomain).PlotGraph") {
			t.Errorf("err not contain (IDomain).PlotGraph: %v", err)
		}
	})

	t.Run("failure(SaveImage)", func(t *testing.T) {
		t.Parallel()
		u := &UseCase{
			repository: &repositoryMock{
				SUMServiceCostGCPFunc: func(ctx context.Context, billingTable string, billingProject string, from time.Time, to time.Time, tz *time.Location, costThreshold float64) ([]domain.GCPServiceCost, error) {
					return tests.NewGCPServiceCosts(tests.TestDate, "test-project", "TestService", 123.45, 1, "USD", 5), nil
				},
				DailyServiceCostGCPFunc: func(ctx context.Context, billingTable string, billingProject string, from time.Time, to time.Time, tz *time.Location, costThreshold float64) ([]domain.GCPServiceCost, error) {
					return tests.NewGCPServiceCosts(tests.TestDate, "test-project", "TestService", 123.45, 1, "USD", 5), nil
				},
				DailyServiceCostGCPMapByServiceFunc: func(servicesOrderBySUMServiceCostGCP []string, dailyServiceCostGCP []domain.GCPServiceCost) map[string][]domain.GCPServiceCost {
					return map[string][]domain.GCPServiceCost{"TestService": tests.NewGCPServiceCosts(tests.TestDate, "test-project", "TestService", 123.45, 1, "USD", 5)}
				},
			},
			domain: &domainMock{
				PlotGraphFunc: func(target io.Writer, ps *domain.PlotGraphParameters) error { return nil },
			},
			infra: &infraMock{
				SaveImageFunc: func(ctx context.Context, image []byte, imageName string, message string) error {
					return testz.ErrTestError
				},
			},
		}
		ctx := context.Background()
		buf := bytes.NewBuffer(nil)
		err := u.PlotDailyServiceCostGCP(ctx, buf, &PlotDailyServiceCostGCPParameters{})
		if !errorz.Contains(err, "(IInfra).SaveImage") {
			t.Errorf("err not contain (IInfra).SaveImage: %v", err)
		}
	})
}
