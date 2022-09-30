package config

import (
	"errors"
	"flag"
	"sync"
	"time"

	"github.com/kunitsuinc/ccc/pkg/constz"
	"github.com/kunitsuinc/util.go/env"
)

var ErrFlagOrEnvIsNotSet = errors.New("flag or environment variable is not set")

// nolint: revive,stylecheck
const (
	DEBUG                = "DEBUG"
	TZ                   = "TZ"
	DAYS                 = "DAYS"
	GOOGLE_CLOUD_PROJECT = "GOOGLE_CLOUD_PROJECT"
	BIGQUERY_DATASET     = "BIGQUERY_DATASET"
	BIGQUERY_TABLE       = "BIGQUERY_TABLE"
	IMAGE_FORMAT         = "IMAGE_FORMAT"
	SLACK_TOKEN          = "SLACK_TOKEN"
	SLACK_CHANNEL        = "SLACK_CHANNEL"
)

type config struct {
	Debug              bool
	TimeZone           *time.Location
	Days               int
	GoogleCloudProject string
	BigQueryDataset    string
	BigQueryTable      string
	ImageFormat        string
	SlackToken         string
	SlackChannel       string
}

// nolint: gochecknoglobals
var (
	cfg   config
	cfgMu sync.Mutex
)

func Load() {
	cfgMu.Lock()
	defer cfgMu.Unlock()

	const empty = ""

	var tz string

	flag.BoolVar(&cfg.Debug, "debug", env.BoolOrDefault(DEBUG, false), "Debug")
	flag.StringVar(&tz, "tz", env.StringOrDefault(TZ, time.UTC.String()), "Time Zone for BigQuery")
	flag.IntVar(&cfg.Days, "days", env.IntOrDefault(DAYS, 30), "Days for BigQuery")
	flag.StringVar(&cfg.ImageFormat, "imgfmt", env.StringOrDefault(IMAGE_FORMAT, "png"), "Image Format")
	flag.StringVar(&cfg.GoogleCloudProject, "project", empty, "Google Cloud Project ID")
	flag.StringVar(&cfg.BigQueryDataset, "dataset", empty, "BigQuery Dataset")
	flag.StringVar(&cfg.BigQueryTable, "table", empty, "BigQuery Table name like: gcp_billing_export_v1_FFFFFF_FFFFFF_FFFFFF")
	flag.StringVar(&cfg.SlackToken, "slack-token", empty, "Slack OAuth Token like: xoxb-999999999999-9999999999999-ZZZZZZZZZZZZZZZZZZZZZZZZ")
	flag.StringVar(&cfg.SlackToken, "slack-channel", empty, "Slack Channel name")
	flag.Parse()

	cfg.TimeZone = constz.TimeZone(tz)

	if cfg.GoogleCloudProject == "" {
		cfg.GoogleCloudProject = env.MustString(GOOGLE_CLOUD_PROJECT)
	}

	if cfg.BigQueryDataset == "" {
		cfg.BigQueryDataset = env.MustString(BIGQUERY_DATASET)
	}

	if cfg.BigQueryTable == "" {
		cfg.BigQueryTable = env.MustString(BIGQUERY_TABLE)
	}

	if cfg.SlackToken == "" {
		cfg.SlackToken = env.MustString(SLACK_TOKEN)
	}

	if cfg.SlackChannel == "" {
		cfg.SlackChannel = env.MustString(SLACK_CHANNEL)
	}
}

func Debug() bool                { return cfg.Debug }
func TimeZone() *time.Location   { return cfg.TimeZone }
func Days() int                  { return cfg.Days }
func ImageFormat() string        { return cfg.ImageFormat }
func GoogleCloudProject() string { return cfg.GoogleCloudProject }
func BigQueryDataset() string    { return cfg.BigQueryDataset }
func BigQueryTable() string      { return cfg.BigQueryTable }
func SlackToken() string         { return cfg.SlackToken }
func SlackChannel() string       { return cfg.SlackChannel }
