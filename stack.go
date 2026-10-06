package errorx

import (
	"fmt"
	"io"
	"runtime"
	"slices"
	"strconv"

	"github.com/lysShub/debug-go"
)

type Stack []uintptr

func (s Stack) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v', 's':
		fs := runtime.CallersFrames(s)
		if state.Flag('+') {
			state.Write([]byte{'\n'})
			for {
				f, more := fs.Next()
				io.WriteString(state, f.File)

				var b [32]byte = [32]byte{0: ':'}
				num := strconv.AppendInt(b[1:1], int64(f.Line), 10) // not possible exceed
				b[len(num)+1] = '\n'

				state.Write(b[:len(num)+2])
				if !more {
					break
				}
			}
		} else {
			f, _ := fs.Next()
			io.WriteString(state, f.File)

			var b [32]byte = [32]byte{0: ':'}
			num := strconv.AppendInt(b[1:1], int64(f.Line), 10) // not possible exceed
			state.Write(b[:len(num)+1])
		}
	default:
	}
}

func NewStack() Stack { return newStack() }

//go:noinline
func newStack() Stack {
	var pcs [64]uintptr
	n := runtime.Callers(3, pcs[:])
	return Stack(slices.Clone(pcs[:n]))
}

func New(msg string) error {
	if debug.Debug() {
		debug.NotEqual(msg, "", "errorx: empty message")
	}
	return WithT(StringErr(msg), newStack())
}
func Errorf(f string, args ...any) error {
	msg := fmt.Sprintf(f, args...)
	if debug.Debug() {
		debug.NotEqual(msg, "", "errorx: empty message")
	}
	return WithT(StringErr(msg), newStack())
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
