package local

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/kunitsucom/ccc/pkg/errors"
	"github.com/kunitsucom/ccc/pkg/log"
	osz "github.com/kunitsucom/util.go/os"
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

func (s *Local) SaveImage(ctx context.Context, image []byte, imageName, message string) error {
	s.imageDir = strings.TrimSuffix(s.imageDir, string(os.PathSeparator))

	if err := osz.CheckDir(s.imageDir); err != nil {
		return errors.Errorf("osz.CheckDir: %w", err)
	}

	imageFilePath := fmt.Sprintf("%s%s%s", s.imageDir, string(os.PathSeparator), imageName)
	log.Debugf("imageFilePath: %s", imageFilePath)

	f, err := os.Create(imageFilePath)
	if err != nil {
		return errors.Errorf("os.Create: %w", err)
	}

	if _, err := f.ReadFrom(bytes.NewReader(image)); err != nil {
		return errors.Errorf("(*os.File).ReadFrom: %w", err)
	}

	if err := f.Sync(); err != nil {
		return errors.Errorf("(*os.File).Sync: %w", err)
	}

	if message != "" {
		log.Infof("%s", message)
	}

	return nil
}
