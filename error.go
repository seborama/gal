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

type ErrUnknownOperator struct {
	msg string
}

func newErrUnknownOperator(msg string) *ErrUnknownOperator {
	return &ErrUnknownOperator{
		msg: msg,
	}
}

func (e ErrUnknownOperator) Error() string {
	return "unknown operator: " + e.msg
}
