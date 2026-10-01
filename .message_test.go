package errorx

import (
	"fmt"
	"net"
	"strings"
	"testing"
)

func Test_Message(t *testing.T) {

	t.Run("base", func(t *testing.T) {
		{
			if WithMessage(nil) != nil {
				t.Fatal("expect nil")
			}
		}
		{
			e := WithMessage(net.ErrClosed)
			if e.Error() != net.ErrClosed.Error() {
				t.Fatalf("expect %q, got %q", net.ErrClosed.Error(), e.Error())
			}
			if Message(e) != net.ErrClosed.Error() {
				t.Fatalf("expect %q, got %q", net.ErrClosed.Error(), Message(e))
			}
			s := fmt.Sprintf("%+v", e)
			if s != net.ErrClosed.Error() {
				t.Fatalf("expect %q, got %q", net.ErrClosed.Error(), s)
			}
		}
		{
			e := WithMessage(nil, "12345")
			if e.Error() != "12345" {
				t.Fatalf("expect %q, got %q", "12345", e.Error())
			}
			if Message(e) != "12345" {
				t.Fatalf("expect %q, got %q", "12345", Message(e))
			}
			s := fmt.Sprintf("%+v", e)
			if s != "12345" {
				t.Fatalf("expect %q, got %q", "12345", s)
			}
		}
		{
			e := WithMessage(net.ErrClosed, "12345")
			exp := fmt.Sprintf("%s:\n %s", "12345", net.ErrClosed.Error())
			if e.Error() != exp {
				t.Fatalf("expect %q, got %q", exp, e.Error())
			}
			if Message(e) != "12345" {
				t.Fatalf("expect %q, got %q", "12345", Message(e))
			}
			s := fmt.Sprintf("%+v", e)
			if s != "12345\n"+net.ErrClosed.Error() {
				t.Fatalf("expect %q, got %q", "12345\n"+net.ErrClosed.Error(), s)
			}
		}
	})

	t.Run("with Stack", func(t *testing.T) {
		{
			if WithStack(WithMessage(nil)) != nil {
				t.Fatal("expect nil")
			}
		}
		{
			e := WithStack(WithMessage(net.ErrClosed))
			if e.Error() != net.ErrClosed.Error() {
				t.Fatalf("expect %q, got %q", net.ErrClosed.Error(), e.Error())
			}
			if Message(e) != net.ErrClosed.Error() {
				t.Fatalf("expect %q, got %q", net.ErrClosed.Error(), Message(e))
			}
			s := fmt.Sprintf("%+v", e)
			if !strings.HasPrefix(s, net.ErrClosed.Error()) {
				t.Fatalf("expect prefix %q, got %q", net.ErrClosed.Error(), s)
			}
		}
		{
			e := WithStack(WithMessage(nil, "12345"))
			if e.Error() != "12345" {
				t.Fatalf("expect %q, got %q", "12345", e.Error())
			}
			if Message(e) != "12345" {
				t.Fatalf("expect %q, got %q", "12345", Message(e))
			}
			s := fmt.Sprintf("%+v", e)
			if !strings.HasPrefix(s, "12345") {
				t.Fatalf("expect prefix %q, got %q", "12345", s)
			}
		}
		{
			e := WithStack(WithMessage(net.ErrClosed, "12345"))
			exp := fmt.Sprintf("%s:\n %s", "12345", net.ErrClosed.Error())
			if e.Error() != exp {
				t.Fatalf("expect %q, got %q", exp, e.Error())
			}
			if Message(e) != "12345" {
				t.Fatalf("expect %q, got %q", "12345", Message(e))
			}
			s := fmt.Sprintf("%+v", e)
			if !strings.HasPrefix(s, "12345\n"+net.ErrClosed.Error()) {
				t.Fatalf("expect prefix %q, got %q", "12345\n"+net.ErrClosed.Error(), s)
			}
		}
	})

	t.Run("with Temporary", func(t *testing.T) {
		{
			if WithTemporary(WithMessage(nil)) != nil {
				t.Fatal("expect nil")
			}
		}
		{
			e := WithTemporary(WithMessage(net.ErrClosed))
			if e.Error() != net.ErrClosed.Error() {
				t.Fatalf("expect %q, got %q", net.ErrClosed.Error(), e.Error())
			}
			if Message(e) != net.ErrClosed.Error() {
				t.Fatalf("expect %q, got %q", net.ErrClosed.Error(), Message(e))
			}
			s := fmt.Sprintf("%+v", e)
			if s != net.ErrClosed.Error() {
				t.Fatalf("expect %q, got %q", net.ErrClosed.Error(), s)
			}
		}
		{
			e := WithTemporary(WithMessage(nil, "12345"))
			if e.Error() != "12345" {
				t.Fatalf("expect %q, got %q", "12345", e.Error())
			}
			if Message(e) != "12345" {
				t.Fatalf("expect %q, got %q", "12345", Message(e))
			}
			s := fmt.Sprintf("%+v", e)
			if s != "12345" {
				t.Fatalf("expect %q, got %q", "12345", s)
			}
		}
		{
			e := WithTemporary(WithMessage(net.ErrClosed, "12345"))
			exp := fmt.Sprintf("%s:\n %s", "12345", net.ErrClosed.Error())
			if e.Error() != exp {
				t.Fatalf("expect %q, got %q", exp, e.Error())
			}
			if Message(e) != "12345" {
				t.Fatalf("expect %q, got %q", "12345", Message(e))
			}
			s := fmt.Sprintf("%+v", e)
			if s != "12345\n"+net.ErrClosed.Error() {
				t.Fatalf("expect %q, got %q", "12345\n"+net.ErrClosed.Error(), s)
			}
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

		e2 := WithMessage(New("abc\n123"))
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
}
