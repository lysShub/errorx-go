package errorx

import (
	"errors"
	"fmt"
)

type temporaryErr struct{ err error }
type temporary interface{ Temporary() bool }

func WithTemporary(err error) error {
	if err == nil {
		return nil
	}
	return &temporaryErr{err: err}
}
func Temporary(err error) bool {
	var t temporary
	if errors.As(err, &t) {
		return t.Temporary()
	} else {
		return false
	}
}
func (t *temporaryErr) Error() string   { return t.err.Error() }
func (t *temporaryErr) Unwrap() error   { return t.err }
func (t *temporaryErr) Temporary() bool { return true }
func (w *temporaryErr) Format(s fmt.State, verb rune) {
	fmt.Fprintf(s, fmt.FormatString(s, verb), w.err)
}
