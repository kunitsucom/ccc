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

// nolint: deadcode,unused
type dailySKUCostGCPParameter struct {
	TimeZone          *time.Location
	GCPBillingTable   string
	GCPBillingProject string
	From              string
	To                string
	CostThreshold     float64
}

// nolint: deadcode,unused,varcheck,gochecknoglobals
var dailySKUCostGCPTemplate = template.Must(template.New("DailySKUCostGCP").Parse(`-- DailySKUCostGCP
SELECT
    FORMAT_DATE('%F', usage_start_time, '{{ .TimeZone }}') AS day,
    service.description AS service,
    CONCAT(service.description, ' ', sku.description) AS sku,
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
    day, service, sku, currency
ORDER BY
    day
ASC
;`))

func (c *BigQuery) DailySKUCostGCP(ctx context.Context, billingTable, billingProject string, from, to time.Time, tz *time.Location, costThreshold float64) ([]domain.GCPSKUCost, error) {
	q, err := buildQuery(dailyServiceCostGCPTemplate, dailySKUCostGCPParameter{
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

	results, err := query[domain.GCPSKUCost](ctx, c.client, q)
	if err != nil {
		return nil, errors.Errorf("query: %w", err)
	}

	return results, nil
}
