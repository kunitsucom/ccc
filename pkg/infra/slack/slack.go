package slack

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"

	"github.com/kunitsuinc/ccc/pkg/errors"
	"github.com/kunitsuinc/ccc/pkg/log"
	"github.com/kunitsuinc/util.go/net/http/httputilz"
)

var (
	ErrSlackAPIError = errors.New("slack api error")
	// nolint: gochecknoglobals
	regexSlackAPIError = regexp.MustCompilePOSIX(`{"ok":false,"error":".+"}`)
)

type Slack struct {
	token   string
	channel string
	client  *http.Client
}

func New(token, channel string, opts ...Option) *Slack {
	s := &Slack{
		token:   token,
		channel: channel,
		client:  new(http.Client),
	}

	for _, opt := range opts {
		s = opt(s)
	}

	return s
}

type Option func(s *Slack) *Slack

func (s *Slack) String() string {
	return "Slack"
}

// nolint: cyclop
func (s *Slack) SaveImage(ctx context.Context, image io.Reader, imageName, message string) error {
	requestBody := &bytes.Buffer{}

	mpw := multipart.NewWriter(requestBody)
	part, err := mpw.CreateFormFile("file", imageName)
	if err != nil {
		return errors.Errorf("(*multipart.Writer).CreateFormFile: %w", err)
	}

	if _, err := io.Copy(part, image); err != nil {
		return errors.Errorf("(io.Writer).Write: %w", err)
	}
	if err := mpw.WriteField("token", s.token); err != nil {
		return errors.Errorf("(*multipart.Writer).WriteField: %w", err)
	}
	if message != "" {
		if err := mpw.WriteField("initial_comment", message); err != nil {
			return errors.Errorf("(*multipart.Writer).WriteField: %w", err)
		}
	}
	if err := mpw.WriteField("channels", s.channel); err != nil {
		return errors.Errorf("(*multipart.Writer).WriteField: %w", err)
	}
	if err := mpw.Close(); err != nil {
		return errors.Errorf("(*multipart.Writer).Close: %w", err)
	}

	requestSlack, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://slack.com/api/files.upload", requestBody)
	if err != nil {
		return errors.Errorf("http.NewRequestWithContext: %w", err)
	}
	requestSlack.Header.Set("content-type", mpw.FormDataContentType())

	responseSlack, err := s.client.Do(requestSlack)
	if err != nil {
		return errors.Errorf("(*http.Client).Do: %w", err)
	}
	defer responseSlack.Body.Close()

	dump, responseBody, err := httputilz.DumpResponse(responseSlack)
	if err != nil {
		return errors.Errorf("httputilz.DumpResponse: %w", err)
	}
	log.Debugf(string(dump))

	if responseSlack.StatusCode >= 300 || regexSlackAPIError.Match(responseBody.Bytes()) {
		return errors.Errorf("%s: %w", strings.ReplaceAll(responseBody.String(), "\n", "\\n"), ErrSlackAPIError)
	}

	return nil
}
