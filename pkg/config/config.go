package config

import (
	"flag"
	"log"
	"sync"
	"time"

	"github.com/kunitsuinc/ccc/pkg/constz"
	"github.com/kunitsuinc/ccc/pkg/errors"
	"github.com/kunitsuinc/util.go/env"
)

var ErrFlagOrEnvIsNotEnough = errors.New("config: flag or environment variable is not enough")

// nolint: revive,stylecheck
const (
	DEBUG                = "DEBUG"
	TZ                   = "TZ"
	DAYS                 = "DAYS"
	GOOGLE_CLOUD_PROJECT = "GOOGLE_CLOUD_PROJECT"
	GCP_BILLING_PROJECT  = "GCP_BILLING_PROJECT"
	GCP_BILLING_TABLE    = "GCP_BILLING_TABLE"
	IMAGE_FORMAT         = "IMAGE_FORMAT"
	MESSAGE              = "MESSAGE"
	SLACK_TOKEN          = "SLACK_TOKEN"
	SLACK_CHANNEL        = "SLACK_CHANNEL"
	IMAGE_DIR            = "IMAGE_DIR"
)

type config struct {
	Debug              bool
	TimeZone           *time.Location
	Days               int
	GoogleCloudProject string
	GCPBillingProject  string
	GCPBillingTable    string
	ImageFormat        string
	Message            string
	SlackToken         string
	SlackChannel       string
	ImageDir           string
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
	flag.StringVar(&cfg.ImageFormat, "image-format", env.StringOrDefault(IMAGE_FORMAT, "png"), "Image Format")
	flag.StringVar(&cfg.GoogleCloudProject, "project", "", "Google Cloud Project ID")
	flag.StringVar(&cfg.GCPBillingTable, "billing-table", "", "GCP Billing export BigQuery Table name like: project-id.dataset_id.gcp_billing_export_v1_FFFFFF_FFFFFF_FFFFFF")
	flag.StringVar(&cfg.GCPBillingProject, "billing-project", "", "Project ID in GCP Billing export BigQuery Table")
	flag.StringVar(&cfg.Message, "message", env.StringOrDefault(MESSAGE, ""), "Slack Message or Log Message")
	flag.StringVar(&cfg.SlackToken, "slack-token", env.StringOrDefault(SLACK_TOKEN, ""), "Slack OAuth Token like: xoxb-999999999999-9999999999999-ZZZZZZZZZZZZZZZZZZZZZZZZ")
	flag.StringVar(&cfg.SlackChannel, "slack-channel", env.StringOrDefault(SLACK_CHANNEL, ""), "Slack Channel name")
	flag.StringVar(&cfg.ImageDir, "image-dir", env.StringOrDefault(IMAGE_DIR, ""), "Directory to save image file")
	flag.Parse()

	cfg.TimeZone = constz.TimeZone(tz)
}

// nolint: cyclop
func Check() error {
	if cfg.GoogleCloudProject == "" {
		v, err := env.String(GOOGLE_CLOUD_PROJECT)
		if err != nil {
			return errors.Errorf("env.String: %w", err)
		}
		cfg.GoogleCloudProject = v
	}

	if cfg.GCPBillingTable == "" {
		v, err := env.String(GCP_BILLING_TABLE)
		if err != nil {
			return errors.Errorf("env.String: %w", err)
		}
		cfg.GCPBillingTable = v
	}

	if cfg.GCPBillingProject == "" {
		v, err := env.String(GCP_BILLING_PROJECT)
		if err != nil {
			return errors.Errorf("env.String: %w", err)
		}
		cfg.GCPBillingProject = v
	}

	switch {
	case cfg.SlackToken != "" && cfg.SlackChannel != "":
		break
	case cfg.ImageDir != "":
		break
	default:
		return errors.Errorf("(%s && %s) || %s: %w", SLACK_TOKEN, SLACK_CHANNEL, IMAGE_DIR, ErrFlagOrEnvIsNotEnough)
	}

	if Debug() {
		log.Printf("[DEBUG] cfg: %#v", cfg)
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
func Message() string            { return cfg.Message }
func SlackToken() string         { return cfg.SlackToken }
func SlackChannel() string       { return cfg.SlackChannel }
func ImageDir() string           { return cfg.ImageDir }
