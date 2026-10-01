package errorx

import (
	"fmt"
	"io"
	"runtime"
	"strconv"

	"github.com/lysShub/debug-go"
)

type stack [32]uintptr
type Stack = *stack

func (t Stack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v', 's':
		fs := runtime.CallersFrames(t[:])
		if s.Flag('+') {
			s.Write([]byte{'\n'})
			for {
				f, more := fs.Next()
				io.WriteString(s, f.File)

				var b [32]byte = [32]byte{0: ':'}
				num := strconv.AppendInt(b[1:1], int64(f.Line), 10) // not possible exceed
				b[len(num)+1] = '\n'

				s.Write(b[:len(num)+2])
				if !more {
					break
				}
			}
		} else {
			f, _ := fs.Next()
			io.WriteString(s, f.File)

			var b [32]byte = [32]byte{0: ':'}
			num := strconv.AppendInt(b[1:1], int64(f.Line), 10) // not possible exceed
			s.Write(b[:len(num)+1])
		}
	default:
	}
}

func NewStack() Stack {
	var s stack
	runtime.Callers(2, s[:])
	return &s
}

//go:noinline
func newStack() Stack {
	var s stack
	runtime.Callers(3, s[:])
	return &s
}

type strerr string             //
func (s strerr) Error() string { return string(s) }

func New(msg string) error {
	if debug.Debug() {
		debug.NotEqual(msg, "", "errorx: empty message")
	}
	return WithT(strerr(msg), newStack())
}
func Errorf(f string, args ...any) error {
	msg := fmt.Sprintf(f, args...)
	if debug.Debug() {
		debug.NotEqual(msg, "", "errorx: empty message")
	}
	return WithT(strerr(msg), newStack())
}

func WithStack(err error) error {
	if err == nil {
		return nil
	}
	stack := newStack()

	if debug.Debug() {
		s := T[Stack](err)
		debug.True(s == nil, "\ncurrent\n", stack, "\nexisted\n", s)
	}
	return WithT(err, stack)
}
