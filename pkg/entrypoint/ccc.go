package entrypoint

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/kunitsuinc/ccc/pkg/config"
	"github.com/kunitsuinc/ccc/pkg/constz"
	"github.com/kunitsuinc/ccc/pkg/errorz"
	"github.com/kunitsuinc/ccc/pkg/infra"
	"github.com/kunitsuinc/ccc/pkg/repository"
	"github.com/kunitsuinc/ccc/pkg/repository/bigquery"
	"github.com/kunitsuinc/ccc/pkg/usecase"
)

func CCC(ctx context.Context) error {
	var (
		tz             = config.TimeZone()
		days           = config.Days()
		projectID      = config.GoogleCloudProject()
		billingTable   = config.GCPBillingTable()
		billingProject = config.GCPBillingProject()
		imageFormat    = config.ImageFormat()
		slackToken     = config.SlackToken()
		slackChannel   = config.SlackChannel()
		slackComment   = config.SlackComment()
	)

	var (
		from = time.Now().In(tz).AddDate(0, 0, -days)
		to   = time.Now().In(tz)
	)

	bq, err := bigquery.New(ctx, projectID)
	if err != nil {
		return errorz.Errorf("bigquery.New: %w", err)
	}

	buf := bytes.NewBuffer(nil)

	r := repository.New(repository.WithBigQuery(bq))
	u := usecase.New(usecase.WithRepository(r))

	if err := u.PlotDailyServiceCostGCP(ctx, buf, billingTable, billingProject, from, to, tz, imageFormat); err != nil {
		return errorz.Errorf("(*usecase.UseCase).PlotDailyServiceCostGCP: %w", err)
	}

	i := infra.New(infra.WithSlack(slackToken))
	if err := i.PostImage(ctx, slackChannel, buf, fmt.Sprintf("%s.%s.%s.%s.%s", projectID, billingProject, billingTable, to.Format(constz.DateOnly), imageFormat), slackComment); err != nil {
		return errorz.Errorf("(*infra.Infra).PostImage: %w", err)
	}

	return nil
}
