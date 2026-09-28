package errorx

import (
	"errors"
	"fmt"
	"log/slog"
	"runtime"
	"strconv"

	"github.com/lysShub/debug-go"
	"go.uber.org/zap/zapcore"
)

type stackErr struct {
	err    error
	stacks *Stacks
}
type stack interface{ Stack() *Stacks }

func WithStack(err error) error {
	if err == nil {
		return nil
	}
	if s := Stack(err); s != nil {
		if debug.Debug() {
			slog.Warn("repeat error stack", slog.Any("current", newStacks()), slog.Any("existed", s))
		}
		return err
	}
	return &stackErr{err: err, stacks: newStacks()}
}

func Stack(err error) *Stacks {
	var s stack
	if errors.As(err, &s) {
		return s.Stack()
	} else {
		return nil
	}
}
func (t *stackErr) Error() string { return t.err.Error() }
func (t *stackErr) Unwrap() error { return t.err }
func (t *stackErr) Stack() *Stacks {
	var s stack
	if errors.As(t.err, &s) {
		return s.Stack() // 仅返回最底层的stack
	}
	return t.stacks
}
func (w *stackErr) Format(s fmt.State, verb rune) {
	fmt.Fprintf(s, fmt.FormatString(s, verb), w.err)
	// stack打印在最后
	if verb == 'v' && s.Flag('+') {
		s.Write([]byte{'\n'})
		w.stacks.Format(s, verb)
	}
}

type Stacks struct {
	n  int32
	pc [64]uintptr
}

var _ slog.LogValuer = (*Stacks)(nil)

func NewStacks(skip ...int) *Stacks { return newStacks(skip...) }

func newStacks(skip ...int) *Stacks {
	s := 3
	if len(skip) > 0 {
		s = skip[0]
	}
	f := &Stacks{}
	n := runtime.Callers(s, f.pc[:])
	n = min(max(0, n-1), len(f.pc))
	f.n = int32(n)
	return f
}
func (f *Stacks) Frames() []uintptr { return f.pc[:f.n] }
func (f *Stacks) LogValue() slog.Value {
	if f == nil {
		return slog.GroupValue()
	}

	var attrs []slog.Attr
	fs := runtime.CallersFrames(f.pc[:f.n])
	for {
		f, more := fs.Next()

		var b = make([]byte, 0, len(f.File)+8)
		b = append(b, f.File...)
		b = append(b, ':')
		b = strconv.AppendInt(b, int64(f.Line), 10)

		attrs = append(attrs, slog.Attr{
			Key:   strconv.Itoa(len(attrs)),
			Value: slog.StringValue(str(b)),
		})
		if !more {
			break
		}
	}
	return slog.GroupValue(attrs...)
}

var _ zapcore.ObjectMarshaler = (*Stacks)(nil)

func (f *Stacks) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if f == nil {
		return nil
	}

	fs := runtime.CallersFrames(f.pc[:f.n])
	for i := 0; ; i++ {
		fr, more := fs.Next()

		var b = make([]byte, 0, len(fr.File)+8)
		b = append(b, fr.File...)
		b = append(b, ':')
		b = strconv.AppendInt(b, int64(fr.Line), 10)
		enc.AddString(strconv.Itoa(i), str(b))

		if !more {
			break
		}
	}
	return nil
}

func (f *Stacks) Format(s fmt.State, verb rune) {
	if verb != 'v' || !s.Flag('+') {
		return
	}

	fs := runtime.CallersFrames(f.pc[:f.n])
	for {
		f, more := fs.Next()
		s.Write(byt(f.File))
		s.Write([]byte{':'})
		s.Write(strconv.AppendInt(make([]byte, 0, 8), int64(f.Line), 10))
		s.Write([]byte{'\n'})
		if !more {
			break
		}
	}
}
