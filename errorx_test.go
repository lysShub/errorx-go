package errorx_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/lysShub/debug-go"
)

func TestMain(m *testing.M) {
	// Assertions abort the process in debug builds; print only so tests can finish.
	debug.Fail = func(s string) {
		fmt.Fprintln(os.Stderr, s)
	}
	os.Exit(m.Run())
}
