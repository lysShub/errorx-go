package errorx

import (
	"fmt"
	"runtime"
	"slices"
	"strconv"
	"unsafe"

	"github.com/lysShub/debug-go"
	"github.com/pkg/errors"
)

// WithStack 为[error]附带堆栈
func WithStack(err error) error {
	return withStack(err)
}

//go:noinline
func withStack(err error) error {
	if err == nil {
		return nil
	}
	var e = &stackErr{
		err: err,
	}
	runtime.Callers(3, e.stacks[:])

	if debug.Debug() {
		var st StackTrace
		debug.False(
			errors.As(err, &st), "current", e.StackTrace(), "existed", st.StackTrace(),
		)
	}
	return e
}

type stackErr struct {
	err    error
	stacks [32]uintptr
}
type StackTrace interface {
	// for compatibility  github.com/pkg/errors
	StackTrace() errors.StackTrace
}

var _ StackTrace = (*stackErr)(nil)

func (t *stackErr) Error() string { return t.err.Error() }
func (t *stackErr) Unwrap() error { return t.err }

func (s *stackErr) StackTrace() errors.StackTrace {
	i := slices.Index(s.stacks[:], uintptr(0))
	if i < 0 && s.stacks[0] != 0 {
		i = len(s.stacks)
	}
	if i > 0 {
		p := (*errors.Frame)(unsafe.Pointer(&s.stacks))
		return unsafe.Slice(p, i)
	}
	return nil
}

func (t *stackErr) Format(s fmt.State, verb rune) {
	if t == nil {
		s.Write(byt("nil"))
		return
	}
	fmt.Fprintf(s, fmt.FormatString(s, verb), t.err)
	if verb == 'v' && s.Flag('+') {
		fs := runtime.CallersFrames(t.stacks[:])
		for {
			f, more := fs.Next()
			s.Write(byt(f.File))

			var b [32]byte = [32]byte{0: ':'}
			num := strconv.AppendInt(b[1:], int64(f.Line), 10) // not possible exceed
			b[len(num)+1] = '\n'
			s.Write(b[:len(num)+2])
			if !more {
				break
			}
		}
	}
}
func byt(s string) []byte { return unsafe.Slice(unsafe.StringData(s), len(s)) }
