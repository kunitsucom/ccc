package slack

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"

	"github.com/kunitsuinc/ccc/pkg/errorz"
	"github.com/kunitsuinc/ccc/pkg/log"
	"github.com/kunitsuinc/util.go/net/http/httputilz"
)

var (
	ErrSlackAPIError = errors.New("slack api error")
	// nolint: gochecknoglobals
	regexSlackAPIError = regexp.MustCompilePOSIX(`{"ok":false,"error":".+"}`)
)

type Slack struct {
	token  string
	client *http.Client
}

func New(token string, opts ...Option) *Slack {
	s := &Slack{
		token:  token,
		client: new(http.Client),
	}

	for _, opt := range opts {
		s = opt(s)
	}

	return s
}

type Option func(s *Slack) *Slack

// nolint: cyclop
func (s *Slack) PostImage(ctx context.Context, slackChannel string, image io.Reader, imageName, comment string) error {
	requestBody := &bytes.Buffer{}

	mpw := multipart.NewWriter(requestBody)
	part, err := mpw.CreateFormFile("file", imageName)
	if err != nil {
		return errorz.Errorf("(*multipart.Writer).CreateFormFile: %w", err)
	}

	if _, err := io.Copy(part, image); err != nil {
		return errorz.Errorf("(io.Writer).Write: %w", err)
	}
	if err := mpw.WriteField("token", s.token); err != nil {
		return errorz.Errorf("(*multipart.Writer).WriteField: %w", err)
	}
	if comment != "" {
		if err := mpw.WriteField("initial_comment", comment); err != nil {
			return errorz.Errorf("(*multipart.Writer).WriteField: %w", err)
		}
	}
	if err := mpw.WriteField("channels", slackChannel); err != nil {
		return errorz.Errorf("(*multipart.Writer).WriteField: %w", err)
	}
	if err := mpw.Close(); err != nil {
		return errorz.Errorf("(*multipart.Writer).Close: %w", err)
	}

	requestSlack, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://slack.com/api/files.upload", requestBody)
	if err != nil {
		return errorz.Errorf("http.NewRequestWithContext: %w", err)
	}
	requestSlack.Header.Set("content-type", mpw.FormDataContentType())

	responseSlack, err := s.client.Do(requestSlack)
	if err != nil {
		return errorz.Errorf("(*http.Client).Do: %w", err)
	}
	defer responseSlack.Body.Close()

	dump, responseBody, err := httputilz.DumpResponse(responseSlack)
	if err != nil {
		return errorz.Errorf("httputilz.DumpResponse: %w", err)
	}
	log.Debugf(string(dump))

	if responseSlack.StatusCode >= 300 || regexSlackAPIError.Match(responseBody.Bytes()) {
		return errorz.Errorf("%s: %w", strings.ReplaceAll(responseBody.String(), "\n", "\\n"), ErrSlackAPIError)
	}

	return nil
}
