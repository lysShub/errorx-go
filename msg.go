package errorx

func WithMessage(err error, msg string) error {
	return WithT(err, msg)
}

func Message(err error) string {
	return T[string](err)
}
