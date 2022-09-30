// nolint: dupl
package bigquery

import (
	"context"
	"text/template"
	"time"

	"github.com/kunitsuinc/ccc/pkg/constz"
	"github.com/kunitsuinc/ccc/pkg/domain"
	"github.com/kunitsuinc/ccc/pkg/errorz"
	"github.com/kunitsuinc/ccc/pkg/log"
)

// nolint: deadcode,unused
type dailySKUCostGCPParameter struct {
	TimeZone      *time.Location
	ProjectID     string
	DatasetName   string
	TableName     string
	From          string
	To            string
	CostThreshold float64
}

// nolint: deadcode,unused,varcheck,gochecknoglobals
var dailySKUCostGCPTemplate = template.Must(template.New("DailySKUCostGCP").Parse(`-- DailySKUCostGCP
SELECT
    FORMAT_DATE('%F', usage_start_time, '{{ .TimeZone }}') AS day,
    IFNULL(project.name, '{{ .ProjectID }}') AS project,
    service.description AS service,
    CONCAT(service.description, ' ', sku.description) AS sku,
    ROUND(SUM(cost * 100)) / 100 AS cost,
    currency
FROM
    ` + "`{{ .ProjectID }}.{{ .DatasetName }}.{{ .TableName }}`" + `
WHERE
    DATE(usage_start_time, '{{ .TimeZone }}') >= DATE("{{ .From }}", '{{ .TimeZone }}')
AND
    DATE(usage_start_time, '{{ .TimeZone }}') <= DATE("{{ .To }}", '{{ .TimeZone }}')
AND
    cost >= {{ .CostThreshold }}
GROUP BY
    day, project, service, sku, currency
ORDER BY
    cost
DESC
;`))

func (c *BigQuery) DailySKUCostGCP(ctx context.Context, projectID, datasetName, tableName string, from, to time.Time, tz *time.Location, costThreshold float64) ([]domain.GCPSKUCost, error) {
	q, err := buildQuery(dailyServiceCostGCPTemplate, dailySKUCostGCPParameter{
		TimeZone:      tz,
		ProjectID:     projectID,
		DatasetName:   datasetName,
		TableName:     tableName,
		From:          from.Format(constz.DateOnly),
		To:            to.Format(constz.DateOnly),
		CostThreshold: costThreshold,
	})
	if err != nil {
		return nil, errorz.Errorf("buildQuery: %w", err)
	}

	log.Debugf(q)

	results, err := query[domain.GCPSKUCost](ctx, c.client, q)
	if err != nil {
		return nil, errorz.Errorf("query: %w", err)
	}

	return results, nil
}
