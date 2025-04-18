// Package errs is a wrapper around standard "errors" package with support for args (fmt.Errorf)
package errs

import (
	"errors"
	"fmt"
)

// ErrUnsupported is proxy for errors.ErrUnsupported
var ErrUnsupported = errors.ErrUnsupported

// New is proxy for errors.New
func New(text string) error {
	return errors.New(text)
}

// Newf is proxy for fmt.Errorf
func Newf(text string, args ...any) error {
	return fmt.Errorf(text, args...)
}

// Errorf is proxy for fmt.Errorf
func Errorf(text string, args ...any) error {
	return fmt.Errorf(text, args...)
}

// Wrap is wraping provided error with prefix using fmt.Errorf and %w
func Wrap(err error, prefix string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf(prefix+": %w", err)
}

// Wrapf is wraping provided error with prefix using fmt.Errorf and %w.
// `prefix` will be expanded with provided args
func Wrapf(err error, prefix string, args ...any) error {
	if err == nil {
		return nil
	}
	args = append(args, err)
	return fmt.Errorf(prefix+": %w", args...)
}

// Unwrap is proxy for errors.Unwrap
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// Is is proxy for errors.Is
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As is proxy for errors.As
func As(err error, target any) bool {
	return errors.As(err, target)
}

// Join is proxy for errors.Join
func Join(errs ...error) error {
	return errors.Join(errs...)
}
