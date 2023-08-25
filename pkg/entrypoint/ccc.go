package entrypoint

import (
	"bytes"
	"context"
	"time"

	"github.com/kunitsucom/ccc/pkg/config"
	"github.com/kunitsucom/ccc/pkg/domain"
	"github.com/kunitsucom/ccc/pkg/errors"
	"github.com/kunitsucom/ccc/pkg/infra"
	"github.com/kunitsucom/ccc/pkg/infra/local"
	"github.com/kunitsucom/ccc/pkg/infra/slack"
	"github.com/kunitsucom/ccc/pkg/repository"
	"github.com/kunitsucom/ccc/pkg/repository/bigquery"
	"github.com/kunitsucom/ccc/pkg/usecase"
)

func CCC(ctx context.Context) error {
	var (
		tz             = config.TimeZone()
		days           = config.Days()
		projectID      = config.GoogleCloudProject()
		billingTable   = config.GCPBillingTable()
		billingProject = config.GCPBillingProject()
		imageFormat    = config.ImageFormat()
		message        = config.Message()
		slackToken     = config.SlackToken()
		slackChannel   = config.SlackChannel()
		imageDir       = config.ImageDir()
	)

	var (
		from = time.Now().In(tz).AddDate(0, 0, -days)
		to   = time.Now().In(tz)
	)

	bq, err := bigquery.New(ctx, projectID)
	if err != nil {
		return errors.Errorf("bigquery.New: %w", err)
	}
	r := repository.New(repository.WithBigQuery(bq))

	d := domain.New()

	var savers []infra.ImageSaver
	if slackToken != "" && slackChannel != "" {
		savers = append(savers, slack.New(slackToken, slackChannel))
	}
	if imageDir != "" {
		savers = append(savers, local.New(imageDir))
	}
	i := infra.New(savers)

	u := usecase.New(usecase.WithRepository(r), usecase.WithDomain(d), usecase.WithInfra(i))

	if err := u.PlotDailyServiceCostGCP(
		ctx,
		bytes.NewBuffer(nil),
		&usecase.PlotDailyServiceCostGCPParameters{
			BillingTable:   billingTable,
			BillingProject: billingProject,
			From:           from,
			To:             to,
			TimeZone:       tz,
			ImageFormat:    imageFormat,
			Message:        message,
		}); err != nil {
		return errors.Errorf("(*usecase.UseCase).PlotDailyServiceCostGCP: %w", err)
	}

	return nil
}
