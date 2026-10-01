package errorx_test

import (
	"fmt"
	"net"
	"strings"
	"testing"

	"github.com/lysShub/errorx-go"
)

func st0() errorx.Stack {
	return errorx.NewStack()
}
func st1() errorx.Stack {
	f := st0()
	if f != nil {
		return f
	}
	return nil
}

func err0() error {
	return errorx.WithStack(net.ErrClosed)
}
func err1() error {
	err := err0()
	if err != nil {
		return err
	}
	return nil
}

var frameFile = "stack_test.go:"
var stackFile = "stack_test.go:"

func Test_New(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		var e = errorx.New("1234")
		if e == nil {
			t.Fatal("expect not nil")
		}
		if e.Error() != "1234" {
			t.Fatalf("expect %q, got %q", "1234", e.Error())
		}
		if errorx.T[errorx.Stack](e) == nil {
			t.Fatal("expect stack not nil")
		}

		s := fmt.Sprintf("%+v", e)
		if !strings.HasPrefix(s, "1234") {
			t.Fatalf("expect prefix %q, got %q", "1234", s)
		}
		if !strings.Contains(s, stackFile) {
			t.Fatalf("expect contains %q, got %q", stackFile, s)
		}
	})

	t.Run("Errorf", func(t *testing.T) {
		var e = errorx.Errorf("%d-%s", 1, "a")
		if e.Error() != "1-a" {
			t.Fatalf("expect %q, got %q", "1-a", e.Error())
		}
		if errorx.T[errorx.Stack](e) == nil {
			t.Fatal("expect stack not nil")
		}
	})

	t.Run("empty message", func(t *testing.T) {
		// 非 debug 构建下不 panic, 返回非 nil 的空消息错误
		var e = errorx.New("")
		if e == nil {
			t.Fatal("expect not nil")
		}
		if e.Error() != "" {
			t.Fatalf("expect empty, got %q", e.Error())
		}
		if errorx.Errorf("%s", "") == nil {
			t.Fatal("expect not nil")
		}
	})

	t.Run("with stack", func(t *testing.T) {
		var e = errorx.WithStack(errorx.New("1234"))
		if e.Error() != "1234" {
			t.Fatalf("expect %q, got %q", "1234", e.Error())
		}
		if errorx.T[errorx.Stack](e) == nil {
			t.Fatal("expect stack not nil")
		}
		if s := fmt.Sprintf("%+v", e); !strings.Contains(s, stackFile) {
			t.Fatalf("expect contains stack, got %q", s)
		}
	})

	t.Run("with stack nil", func(t *testing.T) {
		if errorx.WithStack(nil) != nil {
			t.Fatal("expect nil")
		}
	})

	t.Run("with temporary", func(t *testing.T) {
		var e = errorx.WithTemp(errorx.New("1234"))
		if e.Error() != "1234" {
			t.Fatalf("expect %q, got %q", "1234", e.Error())
		}
		if !errorx.IsTemp(e) {
			t.Fatal("expect temporary")
		}
		if errorx.T[errorx.Stack](e) == nil {
			t.Fatal("expect stack not nil")
		}
	})
}

func Test_Stack(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var s errorx.Stack
		// nil Stack 调用 Format 不应 panic
		_ = fmt.Sprintf("%+v", s)
	})

	t.Run("fmt", func(t *testing.T) {
		s := fmt.Sprintf("%+v", st1())
		if !strings.Contains(s, frameFile) {
			t.Fatalf("expect contains %q, got %q", frameFile, s)
		}
	})
}

func Test_Trace(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		if errorx.WithStack(nil) != nil {
			t.Fatal("expect nil")
		}

		e := err1()
		if errorx.T[errorx.Stack](e) == nil {
			t.Fatal("expect stack not nil")
		}
		if e.Error() != net.ErrClosed.Error() {
			t.Fatalf("expect %q, got %q", net.ErrClosed.Error(), e.Error())
		}

		e2 := errorx.WithT(e, "12345")
		if errorx.T[errorx.Stack](e) == nil {
			t.Fatal("expect stack not nil")
		}
		if !strings.Contains(e2.Error(), net.ErrClosed.Error()) {
			t.Fatalf("expect contains %q, got %q", net.ErrClosed.Error(), e2.Error())
		}

		e3 := errorx.WithTemp(e)
		if errorx.T[errorx.Stack](e) == nil {
			t.Fatal("expect stack not nil")
		}
		if !strings.Contains(e3.Error(), net.ErrClosed.Error()) {
			t.Fatalf("expect contains %q, got %q", net.ErrClosed.Error(), e3.Error())
		}
	})

	t.Run("fmt", func(t *testing.T) {
		s := fmt.Sprintf("%+v", err1())
		if !strings.Contains(s, net.ErrClosed.Error()) {
			t.Fatalf("expect contains %q, got %q", net.ErrClosed.Error(), s)
		}
		if !strings.Contains(s, stackFile) {
			t.Fatalf("expect contains %q, got %q", stackFile, s)
		}
	})

	t.Run("fmt other verb", func(t *testing.T) {
		e1 := errorx.New("abc\n123")
		s1 := fmt.Sprintf("%q", e1)
		if strings.Contains(s1, "\n") {
			t.Fatalf("expect not contains newline, got %q", s1)
		}
		if !strings.Contains(s1, `\n`) {
			t.Fatalf("expect contains %q, got %q", `\n`, s1)
		}

		e2 := errorx.WithStack(errorx.New("abc\n123"))
		s2 := fmt.Sprintf("%q", e2)
		if strings.Contains(s2, "\n") {
			t.Fatalf("expect not contains newline, got %q", s2)
		}
		if !strings.Contains(s2, `\n`) {
			t.Fatalf("expect contains %q, got %q", `\n`, s2)
		}
		if strings.Contains(s2, stackFile) {
			t.Fatalf("expect not contains %q, got %q", stackFile, s2)
		}
	})

	t.Run("fmt with Message", func(t *testing.T) {
		s := fmt.Sprintf("%+v", errorx.WithT(err1(), "123456"))
		for _, exp := range []string{"123456", net.ErrClosed.Error(), stackFile} {
			if !strings.Contains(s, exp) {
				t.Fatalf("expect contains %q, got %q", exp, s)
			}
		}
	})

	t.Run("fmt with Temperory", func(t *testing.T) {
		s := fmt.Sprintf("%+v", errorx.WithTemp(err1()))
		for _, exp := range []string{net.ErrClosed.Error(), stackFile} {
			if !strings.Contains(s, exp) {
				t.Fatalf("expect contains %q, got %q", exp, s)
			}
		}
	})

	t.Run("multi stack", func(t *testing.T) {
		e0 := err1()
		e1 := errorx.WithStack(e0)
		if errorx.T[errorx.Stack](e1) == nil {
			t.Fatal("expect stack not nil")
		}
		if errorx.T[errorx.Stack](e0) == errorx.T[errorx.Stack](e1) {
			t.Fatal("expect distinct stack (each WithStack adds a new one)")
		}
	})

	t.Run("multi stack fmt", func(t *testing.T) {
		s := fmt.Sprintf("%+v", errorx.WithStack(errorx.WithT(errorx.WithT(err1(), "1111"), "222")))
		for _, exp := range []string{"222", "1111", net.ErrClosed.Error(), stackFile} {
			if !strings.Contains(s, exp) {
				t.Fatalf("expect contains %q, got %q", exp, s)
			}
		}
	})
}
