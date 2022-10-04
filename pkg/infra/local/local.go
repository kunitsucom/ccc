package local

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kunitsuinc/ccc/pkg/errorz"
	"github.com/kunitsuinc/ccc/pkg/log"
	"github.com/kunitsuinc/util.go/osz"
)

type Local struct {
	imageDir string
}

func New(imageDir string, opts ...Option) *Local {
	s := &Local{
		imageDir: imageDir,
	}

	for _, opt := range opts {
		s = opt(s)
	}

	return s
}

type Option func(s *Local) *Local

func (s *Local) String() string {
	return "Local"
}

func (s *Local) SaveImage(ctx context.Context, image io.Reader, imageName, message string) error {
	s.imageDir = strings.TrimSuffix(s.imageDir, string(os.PathSeparator))

	if err := osz.CheckDir(s.imageDir); err != nil {
		return errorz.Errorf("osz.CheckDir: %w", err)
	}

	imageFilePath := fmt.Sprintf("%s%s%s", s.imageDir, string(os.PathSeparator), imageName)
	log.Debugf("imageFilePath: %s", imageFilePath)

	f, err := os.Create(imageFilePath)
	if err != nil {
		return errorz.Errorf("os.Create: %w", err)
	}

	if _, err := f.ReadFrom(image); err != nil {
		return errorz.Errorf("(*os.File).ReadFrom: %w", err)
	}

	if err := f.Sync(); err != nil {
		return errorz.Errorf("(*os.File).Sync: %w", err)
	}

	if message != "" {
		log.Infof("%s", message)
	}

	return nil
}
