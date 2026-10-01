package neterr_test

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/lysShub/debug-go"
	"github.com/lysShub/errorx-go"
	"github.com/lysShub/errorx-go/neterr"
)

func TestMain(m *testing.M) {
	// Assertions abort the process in debug builds; print only so tests can finish.
	debug.Fail = func(s string) {
		fmt.Fprintln(os.Stderr, s)
	}
	os.Exit(m.Run())
}

func Test_Base(t *testing.T) {
	var tests = []struct {
		name  string
		check func(error) bool
		err   error
	}{
		{"ConnectReset", neterr.ConnectReset, neterr.ErrConnectReset},
		{"ConnectRefused", neterr.ConnectRefused, neterr.ErrConnectRefused},
		{"ConnectAborted", neterr.ConnectAborted, neterr.ErrConnectAborted},
		{"NetworkUnreach", neterr.NetworkUnreach, neterr.ErrNetworkUnreach},
		{"BuffSize", neterr.BuffSize, neterr.ErrBuffSize},
		{"AddrNotAvail", neterr.AddrNotAvail, neterr.ErrAddrNotAvail},
		{"NetTimeout", neterr.NetTimeout, neterr.ErrNetTimeout},
		{"AddrInuse", neterr.AddrInuse, neterr.ErrAddrInuse},
		{"WouldBlock", neterr.WouldBlock, neterr.ErrWouldBlock},
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
			if tt.check(nil) {
				t.Fatal("expect not matched: nil")
			}
		})
	}
}

func Test_ConnectRefused(t *testing.T) {
	_, err := (&http.Client{Timeout: time.Second}).Get(`http://localhost:12345`)
	if err == nil {
		t.Fatal("expect error")
	}
	if !neterr.ConnectRefused(err) {
		t.Fatalf("expect connect refused, got %v", err)
	}
}

func Test_ErrAddrInuse(t *testing.T) {
	go func() {
		http.ListenAndServe(":12345", nil)
	}()
	time.Sleep(time.Second)

	err := http.ListenAndServe(":12345", nil)
	if !neterr.AddrInuse(err) {
		t.Fatalf("expect addr inuse, got %v", err)
	}
}

func Test_ErrNetTimeout(t *testing.T) {
	var e1 = errorx.WithT(&net.DNSError{IsTimeout: true}, "temp")
	if !neterr.NetTimeout(e1) {
		t.Fatal("expect net timeout")
	}

	var e2 = errorx.WithStack(errorx.New("error"))
	if neterr.NetTimeout(e2) {
		t.Fatal("expect not net timeout")
	}

	var e3 error = nil
	if neterr.NetTimeout(e3) {
		t.Fatal("expect not net timeout")
	}

	var e4 = errorx.WithStack(os.ErrDeadlineExceeded)
	if !neterr.NetTimeout(e4) {
		t.Fatal("expect net timeout")
	}

	var e5 = errorx.WithStack(neterr.ErrNetTimeout)
	if !neterr.NetTimeout(e5) {
		t.Fatal("expect net timeout")
	}
}
