// nolint: dupl
package bigquery

import (
	"context"
	"text/template"
	"time"

	"github.com/kunitsucom/ccc/pkg/consts"
	"github.com/kunitsucom/ccc/pkg/domain"
	"github.com/kunitsucom/ccc/pkg/errors"
	"github.com/kunitsucom/ccc/pkg/log"
)

type dailyProjectCostGCPParameter struct {
	TimeZone          *time.Location
	GCPBillingTable   string
	GCPBillingProject string
	From              string
	To                string
	CostThreshold     float64
}

// nolint: gochecknoglobals
var dailyProjectCostGCPTemplate = template.Must(template.New("DailyProjectCostGCP").Parse(`-- DailyProjectCostGCP
SELECT
    FORMAT_DATE('%F', usage_start_time, '{{ .TimeZone }}') AS day,
    ROUND(SUM(cost * 100)) / 100 AS cost,
    currency
FROM
    ` + "`{{ .GCPBillingTable }}`" + `
WHERE
    project.id = '{{ .GCPBillingProject }}'
AND
    DATE(usage_start_time, '{{ .TimeZone }}') >= DATE("{{ .From }}", '{{ .TimeZone }}')
AND
    DATE(usage_start_time, '{{ .TimeZone }}') < DATE("{{ .To }}", '{{ .TimeZone }}')
AND
    cost >= {{ .CostThreshold }}
GROUP BY
    day, currency
ORDER BY
    day
ASC
;`))

func (c *BigQuery) DailyProjectCostGCP(ctx context.Context, billingTable, billingProject string, from, to time.Time, tz *time.Location, costThreshold float64) ([]domain.GCPCost, error) {
	q, err := buildQuery(dailyProjectCostGCPTemplate, dailyProjectCostGCPParameter{
		TimeZone:          tz,
		GCPBillingTable:   billingTable,
		GCPBillingProject: billingProject,
		From:              from.Format(consts.DateOnly),
		To:                to.Format(consts.DateOnly),
		CostThreshold:     costThreshold,
	})
	if err != nil {
		return nil, errors.Errorf("buildQuery: %w", err)
	}

	log.Debugf("%s", q)

	results, err := query[domain.GCPCost](ctx, c.client, q)
	if err != nil {
		return nil, errors.Errorf("query: %w", err)
	}

	return results, nil
}
