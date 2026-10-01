package errorx

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"testing"
)

func Test_Temporary(t *testing.T) {
	{
		e := WithTemporary(net.ErrClosed)
		if !IsTemporary(e) {
			t.Fatal("expect temporary")
		}
		if e.Error() != net.ErrClosed.Error() {
			t.Fatalf("expect %q, got %q", net.ErrClosed.Error(), e.Error())
		}
	}
	{
		e := WithStack(WithTemporary(net.ErrClosed))
		if !IsTemporary(e) {
			t.Fatal("expect temporary")
		}
		if e.Error() != net.ErrClosed.Error() {
			t.Fatalf("expect %q, got %q", net.ErrClosed.Error(), e.Error())
		}
	}
	{
		e := &net.DNSError{IsTemporary: true}
		if !IsTemporary(e) {
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

		e2 := WithTemporary(errors.New("abc\n123"))
		s2 := fmt.Sprintf("%q", e2)
		if strings.Contains(s2, "\n") {
			t.Fatalf("expect not contains newline, got %q", s2)
		}
		if !strings.Contains(s2, `\n`) {
			t.Fatalf("expect contains %q, got %q", `\n`, s2)
		}
	}

}
