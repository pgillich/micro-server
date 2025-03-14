package logger

import (
	"errors"
	"fmt"
)

type ErrWrap struct {
	err  error
	werr error
}

// Wrap wraps werr with err.
// Error() returns "err: werr".
// Unwrap() returns werr.
func Wrap(err error, werr error) error {
	return &ErrWrap{err, werr}
}

// WrapIf returns nil, if werr is nil
func WrapIf(err error, werr error) error {
	if werr == nil {
		return nil
	}
	return &ErrWrap{err, werr}
}

func (ew *ErrWrap) Error() string {
	return fmt.Sprintf("%s: %s", ew.err, ew.werr)
}

func (ew *ErrWrap) Unwrap() error {
	return ew.werr
}

func (ew *ErrWrap) Is(target error) bool {
	return errors.Is(ew.err, target) ||
		errors.Is(ew.werr, target)
}

func WrapErr(err error) (error, func(error) error) {
	return err, func(werr error) error {
		return Wrap(err, werr)
	}
}
