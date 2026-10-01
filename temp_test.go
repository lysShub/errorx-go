package errorx_test

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"testing"

	"github.com/lysShub/debug-go"
	"github.com/lysShub/errorx-go"
)

func Test_Temporary(t *testing.T) {
	{
		e := errorx.WithTemp(net.ErrClosed)
		if !errorx.IsTemp(e) {
			t.Fatal("expect temporary")
		}
		if e.Error() != net.ErrClosed.Error() {
			t.Fatalf("expect %q, got %q", net.ErrClosed.Error(), e.Error())
		}
	}
	{
		e := errorx.WithStack(errorx.WithTemp(net.ErrClosed))
		if !errorx.IsTemp(e) {
			t.Fatal("expect temporary")
		}
		if e.Error() != net.ErrClosed.Error() {
			t.Fatalf("expect %q, got %q", net.ErrClosed.Error(), e.Error())
		}
	}
	{
		e := &net.DNSError{IsTemporary: true}
		if !errorx.IsTemp(e) {
			t.Fatal("expect temporary")
		}
	}
	{
		e1 := errors.New("abc\n123")
		s1 := fmt.Sprintf("%q", e1)
		if strings.Contains(s1, "\n") {
			t.Fatalf("expect not contains newline, got %q", s1)
		}
		if !strings.Contains(s1, `\n`) {
			t.Fatalf("expect contains %q, got %q", `\n`, s1)
		}

		e2 := errorx.WithTemp(errors.New("abc\n123"))
		s2 := fmt.Sprintf("%q", e2)
		if strings.Contains(s2, "\n") {
			t.Fatalf("expect not contains newline, got %q", s2)
		}
		if !strings.Contains(s2, `\n`) {
			t.Fatalf("expect contains %q, got %q", `\n`, s2)
		}
	}
}

func Test_Warn(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		e := errorx.Warn("warn")
		if e == nil {
			t.Fatal("expect not nil")
		}
		if e.Error() != "warn" {
			t.Fatalf("expect %q, got %q", "warn", e.Error())
		}
		if !errorx.IsTemp(e) {
			t.Fatal("expect temporary")
		}
		// debug 构建附带堆栈, 非 debug 构建不带栈
		if debug.Debug() {
			if errorx.T[errorx.Stack](e) == nil {
				t.Fatal("expect stack in debug")
			}
		} else {
			if errorx.T[errorx.Stack](e) != nil {
				t.Fatal("expect no stack in release")
			}
		}
	})

	t.Run("Warnf", func(t *testing.T) {
		e := errorx.Warnf("%d-%s", 1, "a")
		if e.Error() != "1-a" {
			t.Fatalf("expect %q, got %q", "1-a", e.Error())
		}
		if !errorx.IsTemp(e) {
			t.Fatal("expect temporary")
		}
	})
}
