package errorx

import (
	"errors"
	"fmt"
	"io"
	"unsafe"

	"github.com/lysShub/debug-go"
)

type with[T any] struct {
	err error
	ext T
}
type iwith[T any] interface {
	With() T
}

func WithT[T any](err error, ext T) error {
	if err == nil {
		if debug.Debug() {
			debug.NotEqual(err, nil, "errorx: with nil error")
		}
		return nil
	}
	return with[T]{err: err, ext: ext}
}
func T[T any](err error) (ext T) {
	var i iwith[T]
	if errors.As(err, &i) {
		return i.With()
	}
	return ext
}
func IsT[T comparable](err error, ext T) bool {
	var i iwith[T]
	if errors.As(err, &i) {
		return i.With() == ext
	}
	return false
}

func (t with[T]) Error() string {
	return t.err.Error()
}
func (t with[T]) Unwrap() error { return t.err }
func (t with[T]) Cause() error  { return t.err }
func (t with[T]) With() T       { return t.ext }

func (w with[T]) Format(s fmt.State, verb rune) {
	format := fmt.FormatString(s, verb)

	fmt.Fprintf(s, format, w.err)
	if unsafe.Sizeof(*new(T)) > 0 {
		ext := fmt.Sprintf(format, w.ext)
		if ext != "" {
			s.Write([]byte{':', ' '})
			io.WriteString(s, ext)
		}
	}
}
