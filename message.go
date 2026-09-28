package errorx

import (
	"errors"
	"fmt"
)

type messageErr struct {
	// 至少确保一个是有效值
	err error
	msg string
}
type message interface{ Message() string }

func WithMessage(err error, msg ...string) error {
	switch len(msg) {
	case 0:
		if err == nil {
			return nil
		}
		return &messageErr{err: err, msg: ""}
	case 1:
		if err == nil && msg[0] == "" {
			return nil
		}
		return &messageErr{err: err, msg: msg[0]}
	default:
		panic(len(msg))
	}
}
func Message(err error) string {
	var m message
	if errors.As(err, &m) {
		return m.Message()
	} else {
		return ""
	}
}
func (m *messageErr) Error() string {
	if m.msg != "" && m.err != nil {
		return fmt.Sprintf("%s:\n %s", m.msg, m.err.Error())
	} else if m.msg != "" && m.err == nil {
		return m.msg
	} else if m.msg == "" && m.err != nil {
		return m.err.Error()
	} else {
		panic("")
	}
}
func (t *messageErr) Unwrap() error { return t.err }
func (m *messageErr) Message() string {
	if m.msg == "" {
		return m.err.Error()
	} else {
		return m.msg
	}
}
func (w *messageErr) Format(s fmt.State, verb rune) {
	if w.msg != "" && verb == 'v' && s.Flag('+') {
		s.Write(byt(w.msg))
		if w.err != nil {
			s.Write([]byte{'\n'})
			fmt.Fprintf(s, fmt.FormatString(s, verb), w.err)
		}
	} else {
		fmt.Fprintf(s, fmt.FormatString(s, verb), w.err)
	}
}
