package errorx

import (
	"bytes"
	"fmt"
	"log/slog"
	"net"
	"strings"
	"testing"
)

var _ = 0
var _ = 0
var _ = 0
var _ = 0
var _ = 0
var _ = 0

func s0() *Stacks {
	return newStacks()
}
func s1() *Stacks {
	f := s0() // line 25
	if f != nil {
		return f
	}
	return nil
}

var _ = 0
var _ = 0
var _ = 0

func e0() error {
	var _ = 0
	var _ = 0
	var _ = 0
	err := WithStack(net.ErrClosed) // line 40
	return err
}

func e1() error {
	err := e0()
	if err != nil {
		return err
	}
	return nil
}

var frame1 = "stack_test.go:25"
var stack1 = "stack_test.go:40"

func Test_Stack(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var stack *Stacks

		s := fmt.Sprintf("%+v", stack)
		if s != "<nil>" {
			t.Fatalf("expect %q, got %q", "<nil>", s)
		}

		var b = &bytes.Buffer{}
		slog.New(slog.NewJSONHandler(b, &slog.HandlerOptions{})).Info("-", slog.Any("stack", stack))
		if strings.Contains(b.String(), "stack") {
			t.Fatalf("expect not contains %q, got %q", "stack", b.String())
		}
	})

	t.Run("fmt", func(t *testing.T) {
		s := fmt.Sprintf("%+v", s1())
		ss := strings.Split(s, "\n")
		if !strings.HasSuffix(ss[0], frame1) {
			t.Fatalf("expect suffix %q, got %q", frame1, ss[0])
		}
		if strings.Contains(s, ".s:") {
			t.Fatalf("expect not contains %q, got %q", ".s:", s)
		}
	})

	t.Run("slog", func(t *testing.T) {
		var b = &bytes.Buffer{}
		f := s1()
		slog.New(slog.NewJSONHandler(b, &slog.HandlerOptions{})).Info("-", slog.Any("stack", f))

		s := b.String()
		if !strings.Contains(s, frame1) {
			t.Fatalf("expect contains %q, got %q", frame1, s)
		}
		if strings.Contains(s, ".s:") {
			t.Fatalf("expect not contains %q, got %q", ".s:", s)
		}
	})
}

func Test_Trace(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		if WithStack(nil) != nil {
			t.Fatal("expect nil")
		}

		e1 := e1()
		if Stack(e1) == nil {
			t.Fatal("expect stack not nil")
		}
		if e1.Error() != net.ErrClosed.Error() {
			t.Fatalf("expect %q, got %q", net.ErrClosed.Error(), e1.Error())
		}

		e2 := WithMessage(e1, "12345")
		if Stack(e1) == nil {
			t.Fatal("expect stack not nil")
		}
		if !strings.Contains(e2.Error(), net.ErrClosed.Error()) {
			t.Fatalf("expect contains %q, got %q", net.ErrClosed.Error(), e2.Error())
		}

		e3 := WithTemporary(e1)
		if Stack(e1) == nil {
			t.Fatal("expect stack not nil")
		}
		if !strings.Contains(e3.Error(), net.ErrClosed.Error()) {
			t.Fatalf("expect contains %q, got %q", net.ErrClosed.Error(), e3.Error())
		}
	})

	t.Run("fmt", func(t *testing.T) {
		s := fmt.Sprintf("%+v", e1())

		ss := strings.Split(s, "\n")
		if ss[0] != net.ErrClosed.Error() {
			t.Fatalf("expect %q, got %q", net.ErrClosed.Error(), ss[0])
		}
		if !strings.HasSuffix(ss[1], stack1) {
			t.Fatalf("expect suffix %q, got %q", stack1, ss[1])
		}
		if strings.Contains(s, ".s:") {
			t.Fatalf("expect not contains %q, got %q", ".s:", s)
		}
	})

	t.Run("fmt other verb", func(t *testing.T) {
		e1 := New("abc\n123")
		s1 := fmt.Sprintf("%q", e1)
		if strings.Contains(s1, "\n") {
			t.Fatalf("expect not contains newline, got %q", s1)
		}
		if !strings.Contains(s1, `\n`) {
			t.Fatalf("expect contains %q, got %q", `\n`, s1)
		}

		e2 := WithStack(New("abc\n123"))
		s2 := fmt.Sprintf("%q", e2)
		if strings.Contains(s2, "\n") {
			t.Fatalf("expect not contains newline, got %q", s2)
		}
		if !strings.Contains(s2, `\n`) {
			t.Fatalf("expect contains %q, got %q", `\n`, s2)
		}
		if strings.Contains(s2, stack1) {
			t.Fatalf("expect not contains %q, got %q", stack1, s2)
		}
	})

	t.Run("fmt with Message", func(t *testing.T) {
		s := fmt.Sprintf("%+v", WithMessage(e1(), "123456"))

		ss := strings.Split(s, "\n")
		if ss[0] != "123456" {
			t.Fatalf("expect %q, got %q", "123456", ss[0])
		}
		if ss[1] != net.ErrClosed.Error() {
			t.Fatalf("expect %q, got %q", net.ErrClosed.Error(), ss[1])
		}
		if !strings.HasSuffix(ss[2], stack1) {
			t.Fatalf("expect suffix %q, got %q", stack1, ss[2])
		}
		if strings.Contains(s, ".s:") {
			t.Fatalf("expect not contains %q, got %q", ".s:", s)
		}
	})

	t.Run("fmt with Temperory", func(t *testing.T) {
		s := fmt.Sprintf("%+v", WithTemporary(e1()))

		ss := strings.Split(s, "\n")
		if ss[0] != net.ErrClosed.Error() {
			t.Fatalf("expect %q, got %q", net.ErrClosed.Error(), ss[0])
		}
		if !strings.HasSuffix(ss[1], stack1) {
			t.Fatalf("expect suffix %q, got %q", stack1, ss[1])
		}
		if strings.Contains(s, ".s:") {
			t.Fatalf("expect not contains %q, got %q", ".s:", s)
		}
	})

	t.Run("multi stack", func(t *testing.T) {
		{
			var e0 = e1()
			var e1 = WithStack(e0)
			if Stack(e0) != Stack(e1) {
				t.Fatal("expect equal stack")
			}
		}
		{
			var e0 = e1()
			var e1 = WithStack(WithMessage(e0))
			if Stack(e0) != Stack(e1) {
				t.Fatal("expect equal stack")
			}
		}
		{
			var e0 = e1()
			var e1 = WithStack(WithTemporary(e0))
			if Stack(e0) != Stack(e1) {
				t.Fatal("expect equal stack")
			}
		}
		{
			var e0 = e1()
			var e1 = WithMessage(WithStack(e0))
			if Stack(e0) != Stack(e1) {
				t.Fatal("expect equal stack")
			}
		}
		{
			var e0 = e1()
			var e1 = WithTemporary(WithStack(e0))
			if Stack(e0) != Stack(e1) {
				t.Fatal("expect equal stack")
			}
		}
	})

	t.Run("multi stack fmt", func(t *testing.T) {
		var e = WithStack(WithMessage(WithMessage(e1(), "1111"), "222"))

		s := fmt.Sprintf("%+v", e)

		ss := strings.Split(s, "\n")
		if ss[0] != "222" {
			t.Fatalf("expect %q, got %q", "222", ss[0])
		}
		if ss[1] != "1111" {
			t.Fatalf("expect %q, got %q", "1111", ss[1])
		}
		if ss[2] != net.ErrClosed.Error() {
			t.Fatalf("expect %q, got %q", net.ErrClosed.Error(), ss[2])
		}
		if !strings.HasSuffix(ss[3], stack1) {
			t.Fatalf("expect suffix %q, got %q", stack1, ss[3])
		}
	})
}
