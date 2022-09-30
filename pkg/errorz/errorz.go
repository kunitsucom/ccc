package errorz

import "golang.org/x/xerrors"

// Errorf formats according to a format specifier and returns the string as a
// value that satisfies error.
func Errorf(format string, a ...interface{}) error {
	//nolint: wrapcheck
	return xerrors.Errorf(format, a...)
}
