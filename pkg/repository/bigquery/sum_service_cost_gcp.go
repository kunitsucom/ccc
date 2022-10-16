// nolint: dupl
package bigquery

import (
	"context"
	"text/template"
	"time"

	"github.com/kunitsuinc/ccc/pkg/consts"
	"github.com/kunitsuinc/ccc/pkg/domain"
	"github.com/kunitsuinc/ccc/pkg/errors"
	"github.com/kunitsuinc/ccc/pkg/log"
)

type sumServiceCostGCPParameter struct {
	TimeZone          *time.Location
	GCPBillingTable   string
	GCPBillingProject string
	From              string
	To                string
	CostThreshold     float64
}

// nolint: gochecknoglobals
var sumServiceCostGCPTemplate = template.Must(template.New("SUMServiceCostGCP").Parse(`-- SUMServiceCostGCP
SELECT
    service.description AS service,
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
    service, currency
ORDER BY
    cost
ASC
;`))

func (c *BigQuery) SUMServiceCostGCPAsc(ctx context.Context, billingTable, billingProject string, from, to time.Time, tz *time.Location, costThreshold float64) ([]domain.GCPServiceCost, error) {
	q, err := buildQuery(sumServiceCostGCPTemplate, sumServiceCostGCPParameter{
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

	serviceCostAsc, err := query[domain.GCPServiceCost](ctx, c.client, q)
	if err != nil {
		return nil, errors.Errorf("query: %w", err)
	}

	return serviceCostAsc, nil
}
