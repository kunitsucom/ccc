package domain

type GCPCost struct {
	Day      string  `bigquery:"day"`
	Project  string  `bigquery:"project"`
	Cost     float64 `bigquery:"cost"`
	Currency string  `bigquery:"currency"`
}

type GCPServiceCost struct {
	Day      string  `bigquery:"day"`
	Project  string  `bigquery:"project"`
	Service  string  `bigquery:"service"`
	Cost     float64 `bigquery:"cost"`
	Currency string  `bigquery:"currency"`
}

type GCPSKUCost struct {
	Day      string  `bigquery:"day"`
	Project  string  `bigquery:"project"`
	Service  string  `bigquery:"service"`
	SKU      string  `bigquery:"sku"`
	Cost     float64 `bigquery:"cost"`
	Currency string  `bigquery:"currency"`
}
