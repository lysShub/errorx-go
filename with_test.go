package errorx_test

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/lysShub/errorx-go"
)

type code int

func Test_WithT(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		if errorx.WithT(nil, 1) != nil {
			t.Fatal("expect nil")
		}
	})

	t.Run("value", func(t *testing.T) {
		var base = errors.New("base")
		e := errorx.WithT(base, 5)

		if e.Error() != "base" {
			t.Fatalf("expect %q, got %q", "base", e.Error())
		}
		if errorx.T[int](e) != 5 {
			t.Fatalf("expect %d, got %d", 5, errorx.T[int](e))
		}
		if !errorx.IsT(e, 5) {
			t.Fatal("expect IsT true")
		}
		if errorx.IsT(e, 6) {
			t.Fatal("expect IsT false")
		}
		if !errors.Is(e, base) {
			t.Fatal("expect unwrap to base")
		}
	})

	t.Run("nil error value", func(t *testing.T) {
		if errorx.WithT(nil, net.ErrClosed) != nil {
			t.Fatal("expect nil")
		}
	})

	t.Run("not found", func(t *testing.T) {
		var base = errors.New("base")

		if errorx.T[int](base) != 0 {
			t.Fatalf("expect zero, got %d", errorx.T[int](base))
		}
		if errorx.IsT(base, code(0)) {
			t.Fatal("expect IsT false")
		}
		if errorx.T[errorx.Stack](base) != nil {
			t.Fatal("expect nil stack")
		}
	})

	t.Run("nested", func(t *testing.T) {
		var e = errorx.WithT(errorx.WithT(errors.New("base"), 5), "msg")

		if errorx.T[int](e) != 5 {
			t.Fatalf("expect %d, got %d", 5, errorx.T[int](e))
		}
		if errorx.T[string](e) != "msg" {
			t.Fatalf("expect %q, got %q", "msg", errorx.T[string](e))
		}
	})

	t.Run("format", func(t *testing.T) {
		var e = errorx.WithT(errors.New("base"), "msg")
		if s := fmt.Sprintf("%v", e); s != "base: msg" {
			t.Fatalf("expect %q, got %q", "base: msg", s)
		}

		// an empty string extension is not rendered
		if s := fmt.Sprintf("%v", errorx.WithT(errors.New("base"), "")); s != "base" {
			t.Fatalf("expect %q, got %q", "base", s)
		}
	})

	t.Run("format other verb", func(t *testing.T) {
		var e = errorx.WithT(errors.New("abc\n123"), "msg")
		if s := fmt.Sprintf("%q", e); s != `"abc\n123": "msg"` {
			t.Fatalf("expect %q, got %q", `"abc\n123": "msg"`, s)
		}
	})
}
