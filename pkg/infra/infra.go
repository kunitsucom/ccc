package infra

import (
	"context"
	"io"

	"github.com/kunitsuinc/ccc/pkg/errors"
	"github.com/kunitsuinc/ccc/pkg/log"
)

var (
	ErrImageSaversHaveErrors = errors.New("image savers have errors")
	ErrNoImageSavers         = errors.New("no image savers")
)

type Infra struct {
	imageSavers []ImageSaver
}

type ImageSaver interface {
	String() string
	SaveImage(ctx context.Context, image io.Reader, imageName, message string) error
}

type Option func(i *Infra) *Infra

func New(imageSavers []ImageSaver, opts ...Option) *Infra {
	i := &Infra{
		imageSavers: imageSavers,
	}

	for _, opt := range opts {
		i = opt(i)
	}

	return i
}

func (i *Infra) SaveImage(ctx context.Context, image io.ReadSeeker, imageName, message string) error {
	if len(i.imageSavers) == 0 {
		return ErrNoImageSavers
	}

	var errs []error
	for _, saver := range i.imageSavers {
		if err := saver.SaveImage(ctx, image, imageName, message); err != nil {
			log.Errorf("(ImageSaver).SaveImage: %s: %v", saver, err)
			errs = append(errs, err)
			continue
		}
		if _, err := image.Seek(0, io.SeekStart); err != nil {
			log.Errorf("(io.Seeker).Seek: %v", err)
			errs = append(errs, err)
			continue
		}
	}

	if len(errs) > 0 {
		return errors.Errorf("%v: %w", errs, ErrImageSaversHaveErrors)
	}

	return nil
}
