//go:build unix
// +build unix

package neterr_test

import (
	"errors"
	"testing"

	"github.com/lysShub/errorx-go"
	"github.com/lysShub/errorx-go/neterr"
	"golang.org/x/sys/unix"
)

func Test_BuiltinUnix(t *testing.T) {
	var tests = []struct {
		name  string
		check func(error) bool
		err   error
	}{
		{"NetworkUnreachHost", neterr.NetworkUnreach, unix.EHOSTUNREACH},
		{"BuffSizeNoBufs", neterr.BuffSize, unix.ENOBUFS},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.check(tt.err) {
				t.Fatalf("expect matched: %v", tt.err)
			}
			if !tt.check(errorx.WithStack(tt.err)) {
				t.Fatalf("expect matched through wrapper: %v", tt.err)
			}
			if tt.check(errors.New("other")) {
				t.Fatal("expect not matched: other error")
			}
		})
	}
}
