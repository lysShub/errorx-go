//go:build windows
// +build windows

package neterr

import (
	"errors"
	"net"

	"golang.org/x/sys/windows"
)

var ErrConnectReset = windows.WSAECONNRESET

func connectReset(err error) bool {
	return errors.Is(err, windows.WSAECONNRESET)
}

var ErrConnectRefused = windows.WSAECONNREFUSED

func connectRefused(err error) bool {
	return errors.Is(err, windows.WSAECONNREFUSED) ||
		errors.Is(err, windows.ERROR_CONNECTION_REFUSED)
}

var ErrConnectAborted = windows.WSAECONNABORTED

func connectAborted(err error) bool {
	return errors.Is(err, windows.WSAECONNABORTED)
}

var ErrNetworkUnreach = windows.WSAENETUNREACH

func networkUnreach(err error) bool {
	return errors.Is(err, windows.WSAENETUNREACH) ||
		errors.Is(err, windows.WSAEHOSTUNREACH) ||
		errors.Is(err, windows.ERROR_NETWORK_UNREACHABLE)
}

var ErrBuffSize = windows.ERROR_INSUFFICIENT_BUFFER

func buffSize(err error) bool {
	return errors.Is(err, windows.ERROR_INSUFFICIENT_BUFFER) ||
		errors.Is(err, windows.WSAEMSGSIZE) ||
		errors.Is(err, windows.WSAENOBUFS)
}

var ErrAddrNotAvail = windows.WSAEADDRNOTAVAIL

func addrNotAvail(err error) bool {
	return errors.Is(err, windows.WSAEADDRNOTAVAIL)
}

var ErrNetTimeout = windows.ERROR_TIMEOUT

func netTimeout(err error) bool {
	if e := new(net.Error); errors.As(err, e) && (*e).Timeout() {
		return true
	}
	return errors.Is(err, ErrNetTimeout) || errors.Is(err, windows.WSAETIMEDOUT) || errors.Is(err, windows.WAIT_TIMEOUT)
}

var ErrAddrInuse = windows.WSAEADDRINUSE

func addrInuse(err error) bool {
	return errors.Is(err, ErrAddrInuse)
}

var ErrWouldBlock = windows.WSAEWOULDBLOCK

func wouldBlock(err error) bool {
	return errors.Is(err, windows.WSAEWOULDBLOCK)
}
