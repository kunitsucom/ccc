package config

import (
	"errors"
	"flag"
	"sync"
	"time"

	"github.com/kunitsuinc/ccc/pkg/constz"
	"github.com/kunitsuinc/ccc/pkg/errorz"
	"github.com/kunitsuinc/util.go/env"
)

var ErrFlagOrEnvIsNotSet = errors.New("flag or environment variable is not set")

// nolint: revive,stylecheck
const (
	DEBUG                = "DEBUG"
	TZ                   = "TZ"
	DAYS                 = "DAYS"
	GOOGLE_CLOUD_PROJECT = "GOOGLE_CLOUD_PROJECT"
	GCP_BILLING_PROJECT  = "GCP_BILLING_PROJECT"
	GCP_BILLING_TABLE    = "GCP_BILLING_TABLE"
	IMAGE_FORMAT         = "IMAGE_FORMAT"
	SLACK_TOKEN          = "SLACK_TOKEN"
	SLACK_CHANNEL        = "SLACK_CHANNEL"
	SLACK_COMMENT        = "SLACK_COMMENT"
)

type config struct {
	Debug              bool
	TimeZone           *time.Location
	Days               int
	GoogleCloudProject string
	GCPBillingProject  string
	GCPBillingTable    string
	ImageFormat        string
	SlackToken         string
	SlackChannel       string
	SlackComment       string
}

// nolint: gochecknoglobals
var (
	cfg   config
	cfgMu sync.Mutex
)

func Load() {
	cfgMu.Lock()
	defer cfgMu.Unlock()

	var tz string

	flag.BoolVar(&subcommandVersion, "version", false, "Display version info")
	flag.BoolVar(&cfg.Debug, "debug", env.BoolOrDefault(DEBUG, false), "Debug")
	flag.StringVar(&tz, "tz", env.StringOrDefault(TZ, time.UTC.String()), "Time Zone for BigQuery")
	flag.IntVar(&cfg.Days, "days", env.IntOrDefault(DAYS, 30), "Days for BigQuery")
	flag.StringVar(&cfg.ImageFormat, "imgfmt", env.StringOrDefault(IMAGE_FORMAT, "png"), "Image Format")
	flag.StringVar(&cfg.GoogleCloudProject, "project", "", "Google Cloud Project ID")
	flag.StringVar(&cfg.GCPBillingTable, "billing-table", "", "GCP Billing export BigQuery Table name like: project-id.dataset_id.gcp_billing_export_v1_FFFFFF_FFFFFF_FFFFFF")
	flag.StringVar(&cfg.GCPBillingProject, "billing-project", "", "Project ID in GCP Billing export BigQuery Table")
	flag.StringVar(&cfg.SlackToken, "slack-token", "", "Slack OAuth Token like: xoxb-999999999999-9999999999999-ZZZZZZZZZZZZZZZZZZZZZZZZ")
	flag.StringVar(&cfg.SlackChannel, "slack-channel", "", "Slack Channel name")
	flag.StringVar(&cfg.SlackComment, "slack-comment", env.StringOrDefault(SLACK_COMMENT, ""), "Slack Comment")
	flag.Parse()

	cfg.TimeZone = constz.TimeZone(tz)
}

// nolint: cyclop
func Check() error {
	if cfg.GoogleCloudProject == "" {
		v, err := env.String(GOOGLE_CLOUD_PROJECT)
		if err != nil {
			return errorz.Errorf("env.String: %w", err)
		}
		cfg.GoogleCloudProject = v
	}

	if cfg.GCPBillingTable == "" {
		v, err := env.String(GCP_BILLING_TABLE)
		if err != nil {
			return errorz.Errorf("env.String: %w", err)
		}
		cfg.GCPBillingTable = v
	}

	if cfg.GCPBillingProject == "" {
		v, err := env.String(GCP_BILLING_PROJECT)
		if err != nil {
			return errorz.Errorf("env.String: %w", err)
		}
		cfg.GCPBillingProject = v
	}

	if cfg.SlackToken == "" {
		v, err := env.String(SLACK_TOKEN)
		if err != nil {
			return errorz.Errorf("env.String: %w", err)
		}
		cfg.SlackToken = v
	}

	if cfg.SlackChannel == "" {
		v, err := env.String(SLACK_CHANNEL)
		if err != nil {
			return errorz.Errorf("env.String: %w", err)
		}
		cfg.SlackChannel = v
	}

	return nil
}

func Debug() bool                { return cfg.Debug }
func TimeZone() *time.Location   { return cfg.TimeZone }
func Days() int                  { return cfg.Days }
func ImageFormat() string        { return cfg.ImageFormat }
func GoogleCloudProject() string { return cfg.GoogleCloudProject }
func GCPBillingProject() string  { return cfg.GCPBillingProject }
func GCPBillingTable() string    { return cfg.GCPBillingTable }
func SlackToken() string         { return cfg.SlackToken }
func SlackChannel() string       { return cfg.SlackChannel }
func SlackComment() string       { return cfg.SlackComment }
