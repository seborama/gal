package gal

type ErrSyntaxError struct {
	msg string
}

func newErrSyntaxError(msg string) *ErrSyntaxError {
	return &ErrSyntaxError{
		msg: msg,
	}
}

func (e ErrSyntaxError) Error() string {
	return "syntax error: " + e.msg
}
