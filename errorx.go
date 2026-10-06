// github.com/lysShub/errorx-go

package errorx

type StringErr string

func (s StringErr) Error() string { return string(s) }
