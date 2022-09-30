package infra

import (
	"context"
	"errors"
	"io"

	"github.com/kunitsuinc/ccc/pkg/errorz"
	"github.com/kunitsuinc/ccc/pkg/infra/slack"
)

var ErrNoImagePoster = errors.New("no image poster")

type Infra struct {
	slackToken string
}

type Option func(i *Infra) *Infra

func New(opts ...Option) *Infra {
	i := &Infra{}

	for _, opt := range opts {
		i = opt(i)
	}

	return i
}

func WithSlack(token string) Option {
	return func(i *Infra) *Infra {
		i.slackToken = token
		return i
	}
}

func (i *Infra) PostImage(ctx context.Context, target string, image io.Reader, imageName, comment string) error {
	if i.slackToken != "" {
		s := slack.New(i.slackToken)
		if err := s.PostImage(ctx, target, image, imageName, comment); err != nil {
			return errorz.Errorf("(*slack.Slack).PostImage: %w", err)
		}
		return nil
	}

	return ErrNoImagePoster
}
