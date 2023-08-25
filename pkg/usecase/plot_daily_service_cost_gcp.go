package usecase

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/kunitsucom/ccc/pkg/consts"
	"github.com/kunitsucom/ccc/pkg/domain"
	"github.com/kunitsucom/ccc/pkg/errors"
	"github.com/kunitsucom/ccc/pkg/log"
	slice "github.com/kunitsucom/util.go/slices"
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

func (u *UseCase) PlotDailyServiceCostGCP(ctx context.Context, buf *bytes.Buffer, ps *PlotDailyServiceCostGCPParameters) error {
	sumServiceCostGCPAsc, err := u.repository.SUMServiceCostGCPAsc(ctx, ps.BillingTable, ps.BillingProject, ps.From, ps.To, ps.TimeZone, 0.01)
	if err != nil {
		return errors.Errorf("(IRepository).SUMServiceCostGCP: %w", err)
	}
	servicesOrderBySUMServiceCostAsc := slice.Select(sumServiceCostGCPAsc, func(idx int, source domain.GCPServiceCost) string { return source.Service })

	dailyServiceCostGCP, err := u.repository.DailyServiceCostGCP(ctx, ps.BillingTable, ps.BillingProject, ps.From, ps.To, ps.TimeZone, 0.01)
	log.Debugf("%v", dailyServiceCostGCP)
	if err != nil {
		return errors.Errorf("(IRepository).DailyServiceCostGCP: %w", err)
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
		buf,
		&domain.PlotGraphParameters{
			GraphTitle:        "\n" + fmt.Sprintf("Google Cloud Platform `%s` Cost (from %s to %s)", ps.BillingProject, ps.From.Format(consts.DateOnly), ps.To.Format(consts.DateOnly)),
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
		return errors.Errorf("(IDomain).PlotGraph: %w", err)
	}

	if err := u.infra.SaveImage(ctx, buf.Bytes(), fmt.Sprintf("%s.%s.%s.%s", ps.BillingTable, ps.BillingProject, ps.To.Format(consts.DateOnly), ps.ImageFormat), ps.Message); err != nil {
		return errors.Errorf("(IInfra).SaveImage: %w", err)
	}

	return nil
}
