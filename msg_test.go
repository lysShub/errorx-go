package errorx_test

import (
	"fmt"
	"net"
	"strings"
	"testing"

	"github.com/lysShub/errorx-go"
)

func Test_Message(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		{
			// nil underlying error -> nil
			if errorx.WithMessage(nil, "12345") != nil {
				t.Fatal("expect nil")
			}
		}
		{
			// empty message is not rendered
			e := errorx.WithMessage(net.ErrClosed, "")
			if e.Error() != net.ErrClosed.Error() {
				t.Fatalf("expect %q, got %q", net.ErrClosed.Error(), e.Error())
			}
			if errorx.Message(e) != "" {
				t.Fatalf("expect empty message, got %q", errorx.Message(e))
			}
			if s := fmt.Sprintf("%+v", e); s != net.ErrClosed.Error() {
				t.Fatalf("expect %q, got %q", net.ErrClosed.Error(), s)
			}
		}
		{
			// underlying error plus message
			e := errorx.WithMessage(net.ErrClosed, "12345")
			if e.Error() != net.ErrClosed.Error() {
				t.Fatalf("expect %q, got %q", net.ErrClosed.Error(), e.Error())
			}
			if errorx.Message(e) != "12345" {
				t.Fatalf("expect %q, got %q", "12345", errorx.Message(e))
			}
			exp := net.ErrClosed.Error() + ": 12345"
			if s := fmt.Sprintf("%+v", e); s != exp {
				t.Fatalf("expect %q, got %q", exp, s)
			}
		}
	})

	t.Run("with Stack", func(t *testing.T) {
		{
			if errorx.WithStack(errorx.WithMessage(nil, "")) != nil {
				t.Fatal("expect nil")
			}
		}
		{
			e := errorx.WithStack(errorx.WithMessage(net.ErrClosed, "12345"))
			if e.Error() != net.ErrClosed.Error() {
				t.Fatalf("expect %q, got %q", net.ErrClosed.Error(), e.Error())
			}
			if errorx.Message(e) != "12345" {
				t.Fatalf("expect %q, got %q", "12345", errorx.Message(e))
			}
			s := fmt.Sprintf("%+v", e)
			if !strings.HasPrefix(s, net.ErrClosed.Error()) {
				t.Fatalf("expect prefix %q, got %q", net.ErrClosed.Error(), s)
			}
			if !strings.Contains(s, "12345") {
				t.Fatalf("expect contains %q, got %q", "12345", s)
			}
		}
	})

	t.Run("with Temporary", func(t *testing.T) {
		{
			if errorx.WithTemp(errorx.WithMessage(nil, "")) != nil {
				t.Fatal("expect nil")
			}
		}
		{
			e := errorx.WithTemp(errorx.WithMessage(net.ErrClosed, "12345"))
			if e.Error() != net.ErrClosed.Error() {
				t.Fatalf("expect %q, got %q", net.ErrClosed.Error(), e.Error())
			}
			if errorx.Message(e) != "12345" {
				t.Fatalf("expect %q, got %q", "12345", errorx.Message(e))
			}
			if !errorx.IsTemp(e) {
				t.Fatal("expect temporary")
			}
			exp := net.ErrClosed.Error() + ": 12345"
			if s := fmt.Sprintf("%+v", e); s != exp {
				t.Fatalf("expect %q, got %q", exp, s)
			}
		}
	})

	t.Run("fmt other verb", func(t *testing.T) {
		e1 := errorx.WithMessage(errorx.New("abc\n123"), "msg")
		s1 := fmt.Sprintf("%q", e1)
		if strings.Contains(s1, "\n") {
			t.Fatalf("expect not contains newline, got %q", s1)
		}
		if !strings.Contains(s1, `\n`) {
			t.Fatalf("expect contains %q, got %q", `\n`, s1)
		}
		if strings.Contains(s1, stackFile) {
			t.Fatalf("expect not contains %q, got %q", stackFile, s1)
		}
	})
}
