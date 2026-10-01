//go:build windows
// +build windows

package neterr_test

import (
	"errors"
	"testing"

	"github.com/lysShub/errorx-go"
	"github.com/lysShub/errorx-go/neterr"
	"golang.org/x/sys/windows"
)

func Test_BuiltinWindows(t *testing.T) {
	var tests = []struct {
		name  string
		check func(error) bool
		err   error
	}{
		{"ConnectRefusedError", neterr.ConnectRefused, windows.ERROR_CONNECTION_REFUSED},
		{"NetworkUnreachHost", neterr.NetworkUnreach, windows.WSAEHOSTUNREACH},
		{"NetworkUnreachError", neterr.NetworkUnreach, windows.ERROR_NETWORK_UNREACHABLE},
		{"BuffSizeMsgSize", neterr.BuffSize, windows.WSAEMSGSIZE},
		{"BuffSizeNoBufs", neterr.BuffSize, windows.WSAENOBUFS},
		{"NetTimeoutWsa", neterr.NetTimeout, windows.WSAETIMEDOUT},
		{"NetTimeoutWait", neterr.NetTimeout, windows.WAIT_TIMEOUT},
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
