package errors

import (
	"errors"

	"golang.org/x/xerrors"
)

func New(text string) error {
	// nolint: goerr113
	return errors.New(text)
}

// Errorf formats according to a format specifier and returns the string as a
// value that satisfies error.
func Errorf(format string, a ...interface{}) error {
	//nolint: wrapcheck
	return xerrors.Errorf(format, a...)
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}
