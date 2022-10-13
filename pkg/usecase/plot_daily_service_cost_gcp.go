package usecase

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/kunitsuinc/ccc/pkg/constz"
	"github.com/kunitsuinc/ccc/pkg/domain"
	"github.com/kunitsuinc/ccc/pkg/errors"
	"github.com/kunitsuinc/ccc/pkg/log"
	"github.com/kunitsuinc/util.go/bytez"
	"github.com/kunitsuinc/util.go/slice"
	"gonum.org/v1/plot/plotter"
)

type PlotDailyServiceCostGCPParameters struct {
	BillingTable   string
	BillingProject string
	From           time.Time
	To             time.Time
	TimeZone       *time.Location
	ImageFormat    string
	Message        string
}

func (u *UseCase) PlotDailyServiceCostGCP(ctx context.Context, target io.ReadWriter, ps *PlotDailyServiceCostGCPParameters) error {
	sumServiceCostGCPAsc, err := u.repository.SUMServiceCostGCPAsc(ctx, ps.BillingTable, ps.BillingProject, ps.From, ps.To, ps.TimeZone, 0.01)
	if err != nil {
		return errors.Errorf("(RepositoryIF).SUMServiceCostGCP: %w", err)
	}
	servicesOrderBySUMServiceCostAsc := slice.Select(sumServiceCostGCPAsc, func(idx int, source domain.GCPServiceCost) string { return source.Service })

	dailyServiceCostGCP, err := u.repository.DailyServiceCostGCP(ctx, ps.BillingTable, ps.BillingProject, ps.From, ps.To, ps.TimeZone, 0.01)
	log.Debugf("%v", dailyServiceCostGCP)
	if err != nil {
		return errors.Errorf("(RepositoryIF).DailyServiceCostGCP: %w", err)
	}
	currencies := slice.Uniq(slice.Select(dailyServiceCostGCP, func(_ int, s domain.GCPServiceCost) (selected string) { return s.Currency }))
	if len(currencies) != 1 {
		return errors.Errorf("%s: %s: %v: %w", ps.BillingTable, ps.BillingProject, currencies, ErrMixedCurrenciesDataSourceIsNotSupported)
	}
	currency := currencies[0]
	dailyServiceCostGCPMapByService := u.repository.DailyServiceCostGCPMapByService(servicesOrderBySUMServiceCostAsc, dailyServiceCostGCP)

	dailyServiceCostsForPlot := make(map[string]plotter.Values)
	var xAxisPointsCount int // NOTE: X 軸の数値の数を数える
	for k, v := range dailyServiceCostGCPMapByService {
		dailyServiceCostsForPlot[k] = slice.Select(v, func(_ int, source domain.GCPServiceCost) float64 { return source.Cost })

		log.Debugf("%s: data count: %d", k, len(v))
		if len(v) > xAxisPointsCount {
			xAxisPointsCount = len(v)
		}
	}

	if err := u.domain.PlotGraph(
		target,
		&domain.PlotGraphParameters{
			GraphTitle:        "\n" + fmt.Sprintf("Google Cloud Platform `%s` Cost (from %s to %s)", ps.BillingProject, ps.From.Format(constz.DateOnly), ps.To.Format(constz.DateOnly)),
			XLabelText:        "\n" + fmt.Sprintf("Date (%s)", ps.TimeZone.String()),
			YLabelText:        "\n" + currency,
			Width:             1280,
			Hight:             720,
			XAxisPointsCount:  xAxisPointsCount,
			From:              ps.From,
			To:                ps.To,
			TimeZone:          ps.TimeZone,
			OrderedLegendsAsc: servicesOrderBySUMServiceCostAsc,
			LegendValuesMap:   dailyServiceCostsForPlot,
			ImageFormat:       ps.ImageFormat,
		},
	); err != nil {
		return errors.Errorf("(DomainIF).PlotGraph: %w", err)
	}

	if err := u.infra.SaveImage(ctx, bytez.NewReadSeekBuffer(target), fmt.Sprintf("%s.%s.%s.%s", ps.BillingTable, ps.BillingProject, ps.To.Format(constz.DateOnly), ps.ImageFormat), ps.Message); err != nil {
		return errors.Errorf("(InfraIF).SaveImage: %w", err)
	}

	return nil
}
