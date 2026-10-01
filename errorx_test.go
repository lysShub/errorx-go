package errorx_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/lysShub/debug-go"
)

func TestMain(m *testing.M) {
	// debug 构建下断言默认会终止进程, 测试中改为仅打印, 以便 UT 能跑完
	debug.Fail = func(s string) {
		fmt.Fprintln(os.Stderr, s)
	}
	os.Exit(m.Run())
}
