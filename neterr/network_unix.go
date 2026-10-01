//go:build unix
// +build unix

package neterr

import (
	"errors"
	"net"

	"golang.org/x/sys/unix"
)

var ErrConnectReset = unix.ECONNRESET

func connectReset(err error) bool {
	return errors.Is(err, unix.ECONNRESET)
}

var ErrConnectRefused = unix.ECONNREFUSED

func connectRefused(err error) bool {
	return errors.Is(err, unix.ECONNREFUSED)
}

var ErrConnectAborted = unix.ECONNABORTED

func connectAborted(err error) bool {
	return errors.Is(err, unix.ECONNABORTED)
}

var ErrNetworkUnreach = unix.ENETUNREACH

func networkUnreach(err error) bool {
	return errors.Is(err, unix.ENETUNREACH) || errors.Is(err, unix.EHOSTUNREACH)
}

var ErrBuffSize = unix.EMSGSIZE

func buffSize(err error) bool {
	return errors.Is(err, unix.EMSGSIZE) ||
		errors.Is(err, unix.ENOBUFS)
}

var ErrAddrNotAvail = unix.EADDRNOTAVAIL

func addrNotAvail(err error) bool {
	return errors.Is(err, unix.EADDRNOTAVAIL)
}

var ErrNetTimeout = unix.ETIMEDOUT

func netTimeout(err error) bool {
	if e := new(net.Error); errors.As(err, e) && (*e).Timeout() {
		return true
	}
	return errors.Is(err, unix.ETIMEDOUT)
}

var ErrAddrInuse = unix.EADDRINUSE

func addrInuse(err error) bool {
	return errors.Is(err, ErrAddrInuse)
}

var ErrWouldBlock = unix.EWOULDBLOCK

func wouldBlock(err error) bool {
	return errors.Is(err, unix.EWOULDBLOCK)
}
