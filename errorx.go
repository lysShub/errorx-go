// ref github.com/pkg/errors
package errorx

import (
	"fmt"
	"unsafe"
)

type fundamental struct {
	msg    string
	frames *Stacks
}

func New(s string) error {
	if s == "" {
		panic(s)
	}
	return &fundamental{msg: s, frames: newStacks()}
}
func Errorf(f string, args ...any) error {
	msg := fmt.Sprintf(f, args...)
	if msg == "" {
		panic(msg)
	}
	return &fundamental{
		msg:    msg,
		frames: newStacks(),
	}
}
func (b *fundamental) Error() string  { return b.msg }
func (t *fundamental) Stack() *Stacks { return t.frames }
func (b *fundamental) Format(s fmt.State, verb rune) {
	fmt.Fprintf(s, fmt.FormatString(s, verb), b.msg)
	// stack打印在最后
	if verb == 'v' && s.Flag('+') {
		s.Write([]byte{'\n'})
		b.frames.Format(s, verb)
	}
}

func byt(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
func str(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
