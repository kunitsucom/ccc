// nolint: testpackage
package usecase

import (
	"bytes"
	"context"
	"testing"

	"github.com/kunitsucom/ccc/pkg/domain"
	"github.com/kunitsucom/ccc/pkg/errors"
	"github.com/kunitsucom/ccc/pkg/tests"
errorz	"github.com/kunitsucom/util.go/errors"
testz	"github.com/kunitsucom/util.go/test"
)

func TestUsecase_PlotDailyServiceCostGCP(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		u := &UseCase{
			repository: &repositoryMock{
				SUMServiceCostGCP_GCPServiceCost:                          tests.NewGCPServiceCosts(tests.TestDate, "test-project", "TestService", 123.45, 1, "USD", 5),
				SUMServiceCostGCP_error:                                   nil,
				DailyServiceCostGCP_GCPServiceCost:                        tests.NewGCPServiceCosts(tests.TestDate, "test-project", "TestService", 123.45, 1, "USD", 5),
				DailyServiceCostGCP_error:                                 nil,
				DailyServiceCostGCPMapByService_map_string_GCPServiceCost: map[string][]domain.GCPServiceCost{"TestService": tests.NewGCPServiceCosts(tests.TestDate, "test-project", "TestService", 123.45, 1, "USD", 5)},
			},
			domain: &domainMock{
				PlotGraph_error: nil,
			},
			infra: &infraMock{
				SaveImage_error: nil,
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
				SUMServiceCostGCP_GCPServiceCost: []domain.GCPServiceCost{},
				SUMServiceCostGCP_error:          testz.ErrTestError,
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
				SUMServiceCostGCP_GCPServiceCost: tests.NewGCPServiceCosts(tests.TestDate, "test-project", "TestService", 123.45, 1, "USD", 5),
				SUMServiceCostGCP_error:          nil,
				DailyServiceCostGCP_error:        testz.ErrTestError,
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
				SUMServiceCostGCP_GCPServiceCost:                          append(tests.NewGCPServiceCosts(tests.TestDate, "test-project", "TestService", 123.45, 1, "USD", 5), tests.NewGCPServiceCosts(tests.TestDate, "test-project", "TestService", 123.45, 1, "JPY", 5)...),
				SUMServiceCostGCP_error:                                   nil,
				DailyServiceCostGCP_GCPServiceCost:                        append(tests.NewGCPServiceCosts(tests.TestDate, "test-project", "TestService", 123.45, 1, "USD", 5), tests.NewGCPServiceCosts(tests.TestDate, "test-project", "TestService", 123.45, 1, "JPY", 5)...),
				DailyServiceCostGCP_error:                                 nil,
				DailyServiceCostGCPMapByService_map_string_GCPServiceCost: map[string][]domain.GCPServiceCost{"TestService": tests.NewGCPServiceCosts(tests.TestDate, "test-project", "TestService", 123.45, 1, "USD", 5)},
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
				SUMServiceCostGCP_GCPServiceCost:                          tests.NewGCPServiceCosts(tests.TestDate, "test-project", "TestService", 123.45, 1, "USD", 5),
				SUMServiceCostGCP_error:                                   nil,
				DailyServiceCostGCP_GCPServiceCost:                        tests.NewGCPServiceCosts(tests.TestDate, "test-project", "TestService", 123.45, 1, "USD", 5),
				DailyServiceCostGCP_error:                                 nil,
				DailyServiceCostGCPMapByService_map_string_GCPServiceCost: map[string][]domain.GCPServiceCost{"TestService": tests.NewGCPServiceCosts(tests.TestDate, "test-project", "TestService", 123.45, 1, "USD", 5)},
			},
			domain: &domainMock{
				PlotGraph_error: testz.ErrTestError,
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
				SUMServiceCostGCP_GCPServiceCost:                          tests.NewGCPServiceCosts(tests.TestDate, "test-project", "TestService", 123.45, 1, "USD", 5),
				SUMServiceCostGCP_error:                                   nil,
				DailyServiceCostGCP_GCPServiceCost:                        tests.NewGCPServiceCosts(tests.TestDate, "test-project", "TestService", 123.45, 1, "USD", 5),
				DailyServiceCostGCP_error:                                 nil,
				DailyServiceCostGCPMapByService_map_string_GCPServiceCost: map[string][]domain.GCPServiceCost{"TestService": tests.NewGCPServiceCosts(tests.TestDate, "test-project", "TestService", 123.45, 1, "USD", 5)},
			},
			domain: &domainMock{
				PlotGraph_error: nil,
			},
			infra: &infraMock{
				SaveImage_error: testz.ErrTestError,
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
