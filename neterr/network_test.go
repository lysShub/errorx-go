package neterr_test

import (
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
	// debug 构建下断言默认会终止进程, 测试中改为仅打印, 以便 UT 能跑完
	debug.Fail = func(s string) {
		fmt.Fprintln(os.Stderr, s)
	}
	os.Exit(m.Run())
}

/*
	type timeoutErr struct{ error }

	func Timeout(err error) error {
		if err == nil {
			return nil
		}
		return &timeoutErr{error: err}
	}
	func IsTimeout(err error) bool {
		timeout := UnwrapTo[interface{ Timeout() bool }](err)
		if timeout == nil {
			return false
		} else {
			return timeout.Timeout()
		}
	}
	func (t *timeoutErr) Error() string { return t.error.Error() }
	func (t *timeoutErr) Unwrap() error { return t.error }
	func (t *timeoutErr) Timeout() bool { return true }
*/

func Test_Temporary(t *testing.T) {
	var e1 = errorx.WithT(&net.DNSError{IsTemporary: true}, "temp")
	if !errorx.IsTemp(e1) {
		t.Fatal("expect temporary")
	}

	var e2 = errorx.WithStack(errorx.New("error"))
	if errorx.IsTemp(e2) {
		t.Fatal("expect not temporary")
	}

	var e3 error = nil
	if errorx.IsTemp(e3) {
		t.Fatal("expect not temporary")
	}
}

// func Test_Timeout(t *testing.T) {
// 	var e1 = errorx.WithMessage(&net.DNSError{IsTimeout: true}, "temp")
// 	require.True(t, IsTimeout(e1))
//
// 	var e2 = errorx.WithStack(errorx.New("error"))
// 	require.False(t, IsTimeout(e2))
//
// 	var e3 error = nil
// 	require.False(t, IsTimeout(e3))
// }

func Test_ConnectRefused(t *testing.T) {
	_, err := (&http.Client{Timeout: time.Second}).Get(`http://localhost:12345`)
	if err == nil {
		t.Fatal("expect error")
	}
	if !neterr.ConnectRefused(err) {
		t.Fatalf("expect connect refused, got %v", err)
	}
}

func Test_Builtin(t *testing.T) {
	if !neterr.ConnectRefused(neterr.ErrConnectRefused) {
		t.Fatal("expect connect refused")
	}
	if !neterr.NetworkUnreach(neterr.ErrNetworkUnreach) {
		t.Fatal("expect network unreach")
	}
	if !neterr.BuffSize(neterr.ErrBuffSize) {
		t.Fatal("expect buff size")
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
