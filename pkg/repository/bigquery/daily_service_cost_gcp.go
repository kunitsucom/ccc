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

type dailyServiceCostGCPParameter struct {
	TimeZone      *time.Location
	ProjectID     string
	DatasetName   string
	TableName     string
	From          string
	To            string
	CostThreshold float64
}

// nolint: gochecknoglobals
var dailyServiceCostGCPTemplate = template.Must(template.New("DailyServiceCostGCP").Parse(`-- DailyServiceCostGCP
SELECT
    FORMAT_DATE('%F', usage_start_time, '{{ .TimeZone }}') AS day,
    IFNULL(project.name, '{{ .ProjectID }}') AS project,
    service.description AS service,
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
    day, project, service, currency
ORDER BY
    cost
DESC
;`))

func (c *BigQuery) DailyServiceCostGCP(ctx context.Context, projectID, datasetName, tableName string, from, to time.Time, tz *time.Location, costThreshold float64) ([]domain.GCPServiceCost, error) {
	q, err := buildQuery(dailyServiceCostGCPTemplate, dailyServiceCostGCPParameter{
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

	results, err := query[domain.GCPServiceCost](ctx, c.client, q)
	if err != nil {
		return nil, errorz.Errorf("query: %w", err)
	}

	return results, nil
}
