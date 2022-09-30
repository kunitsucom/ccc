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

type sumServiceCostGCPParameter struct {
	TimeZone      *time.Location
	ProjectID     string
	DatasetName   string
	TableName     string
	From          string
	To            string
	CostThreshold float64
}

// nolint: gochecknoglobals
var sumServiceCostGCPTemplate = template.Must(template.New("SUMServiceCostGCP").Parse(`-- SUMServiceCostGCP
SELECT
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
    project, service, currency
ORDER BY
    cost
DESC
;`))

func (c *BigQuery) SUMServiceCostGCP(ctx context.Context, projectID, datasetName, tableName string, from, to time.Time, tz *time.Location, costThreshold float64) ([]domain.GCPServiceCost, error) {
	q, err := buildQuery(sumServiceCostGCPTemplate, sumServiceCostGCPParameter{
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
