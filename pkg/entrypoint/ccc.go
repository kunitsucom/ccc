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
	"github.com/kunitsuinc/ccc/pkg/infra/local"
	"github.com/kunitsuinc/ccc/pkg/infra/slack"
	"github.com/kunitsuinc/ccc/pkg/repository"
	"github.com/kunitsuinc/ccc/pkg/repository/bigquery"
	"github.com/kunitsuinc/ccc/pkg/usecase"
	"github.com/kunitsuinc/util.go/bytez"
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
		return errorz.Errorf("bigquery.New: %w", err)
	}

	buf := bytes.NewBuffer(nil)

	r := repository.New(repository.WithBigQuery(bq))
	u := usecase.New(usecase.WithRepository(r))

	if err := u.PlotDailyServiceCostGCP(ctx, buf, billingTable, billingProject, from, to, tz, imageFormat); err != nil {
		return errorz.Errorf("(*usecase.UseCase).PlotDailyServiceCostGCP: %w", err)
	}

	var savers []infra.ImageSaver
	if slackToken != "" && slackChannel != "" {
		savers = append(savers, slack.New(slackToken, slackChannel))
	}
	if imageDir != "" {
		savers = append(savers, local.New(imageDir))
	}
	i := infra.New(savers)
	if err := i.SaveImage(ctx, bytez.NewReadSeekBuffer(buf), fmt.Sprintf("%s.%s.%s.%s", billingTable, billingProject, to.Format(constz.DateOnly), imageFormat), message); err != nil {
		return errorz.Errorf("(*infra.Infra).SaveImage: %w", err)
	}

	return nil
}
