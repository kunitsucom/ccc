package bigquery

import (
	"bytes"
	"context"
	"sync"
	"text/template"

	"cloud.google.com/go/bigquery"
	"github.com/kunitsucom/ccc/pkg/errors"
	"google.golang.org/api/iterator"
)

type BigQuery struct {
	client    *bigquery.Client
	projectID string
}

// nolint: gochecknoglobals
var (
	bigqueryClients   = make(map[string]*bigquery.Client)
	bigqueryClientsMu sync.Mutex
)

func New(ctx context.Context, projectID string) (*BigQuery, error) {
	bigqueryClientsMu.Lock()
	defer bigqueryClientsMu.Unlock()

	client := bigqueryClients[projectID]
	if client == nil {
		var err error
		client, err = bigquery.NewClient(ctx, projectID)
		if err != nil {
			return nil, errors.Errorf("bigquery.NewClient: %w", err)
		}

		bigqueryClients[projectID] = client
	}

	return &BigQuery{
		client:    client,
		projectID: projectID,
	}, nil
}

func (c *BigQuery) Renew(ctx context.Context) error {
	bigqueryClientsMu.Lock()
	defer bigqueryClientsMu.Unlock()

	client := bigqueryClients[c.projectID]
	if client == nil {
		var err error
		client, err = bigquery.NewClient(ctx, c.projectID)
		if err != nil {
			return errors.Errorf("bigquery.NewClient: %w", err)
		}

		bigqueryClients[c.projectID] = client
	}

	c.client = client

	return nil
}

func buildQuery(tmpl *template.Template, tmplParams any) (string, error) {
	buf := bytes.NewBuffer(nil)

	if err := tmpl.Execute(buf, tmplParams); err != nil {
		return "", errors.Errorf("(*template.Template).Execute: %w", err)
	}

	return buf.String(), nil
}

func query[Result any](ctx context.Context, c *bigquery.Client, q string) ([]Result, error) {
	job, err := c.Query(q).Run(ctx)
	if err != nil {
		return nil, errors.Errorf("(*bigquery.Query).Run: %w", err)
	}

	ri, err := job.Read(ctx)
	if err != nil {
		return nil, errors.Errorf("(*bigquery.Job).Read: %w", err)
	}

	results := make([]Result, 0)
	for {
		var result Result
		if err := ri.Next(&result); err != nil {
			if errors.Is(err, iterator.Done) {
				break
			}
			return nil, errors.Errorf("(*bigquery.RowIterator).Next: %w", err)
		}

		results = append(results, result)
	}

	return results, nil
}
