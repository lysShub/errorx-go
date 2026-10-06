package errorx

import (
	"errors"
	"fmt"

	"github.com/lysShub/debug-go"
)

type temporary struct{} //
var temp = temporary{}

// Warn creates a temporary-marked error, carrying a stack in debug builds.
func Warn(msg string) error {
	if debug.Debug() {
		debug.NotEqual(msg, "", "errorx: empty message")

		return WithT(WithT(StringErr(msg), newStack()), temp)
	}
	return WithT(StringErr(msg), temp)
}
func Warnf(f string, args ...any) error {
	msg := fmt.Sprintf(f, args...)
	if debug.Debug() {
		debug.NotEqual(msg, "", "errorx: empty message")

		return WithT(WithT(StringErr(msg), newStack()), temp)
	}
	return WithT(StringErr(msg), temp)
}

func WithTemp(err error) error {
	return WithT(err, temp)
}
func IsTemp(err error) bool {
	if IsT(err, temp) {
		return true
	}
	type iface interface{ Temporary() bool }

	var i iface
	if errors.As(err, &i) {
		return i.Temporary()
	}
	return false
}
