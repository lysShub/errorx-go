package errorx

import (
	"fmt"
	"strings"
	"testing"
)

func Test_New(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		func() {
			defer func() {
				if recover() == nil {
					t.Fatal("expect panic")
				}
			}()
			New("")
		}()
		func() {
			defer func() {
				if recover() == nil {
					t.Fatal("expect panic")
				}
			}()
			Errorf("%s", "")
		}()

		var e = New("1234")
		if e.Error() != "1234" {
			t.Fatalf("expect %q, got %q", "1234", e.Error())
		}
		s := fmt.Sprintf("%+v", e)
		ss := strings.Split(s, "\n")
		if ss[0] != "1234" {
			t.Fatalf("expect %q, got %q", "1234", ss[0])
		}
		if !strings.Contains(s, "errorx_test.go:") {
			t.Fatalf("expect contains %q, got %q", "errorx_test.go:", s)
		}
	})

	t.Run("with stack", func(t *testing.T) {
		var e = WithStack(New("1234"))
		if e.Error() != "1234" {
			t.Fatalf("expect %q, got %q", "1234", e.Error())
		}
		s := fmt.Sprintf("%+v", e)
		ss := strings.Split(s, "\n")
		if ss[0] != "1234" {
			t.Fatalf("expect %q, got %q", "1234", ss[0])
		}
		if !strings.Contains(s, "errorx_test.go:") {
			t.Fatalf("expect contains %q, got %q", "errorx_test.go:", s)
		}
	})

	t.Run("with temporary", func(t *testing.T) {
		var e = WithTemporary(New("1234"))
		if e.Error() != "1234" {
			t.Fatalf("expect %q, got %q", "1234", e.Error())
		}
		if !Temporary(e) {
			t.Fatal("expect temporary")
		}

		s := fmt.Sprintf("%+v", e)
		ss := strings.Split(s, "\n")
		if ss[0] != "1234" {
			t.Fatalf("expect %q, got %q", "1234", ss[0])
		}
		if !strings.Contains(s, "errorx_test.go:") {
			t.Fatalf("expect contains %q, got %q", "errorx_test.go:", s)
		}
	})

}
