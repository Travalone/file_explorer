package common

type Error struct {
	msg string
}

func (err *Error) Error() string {
	return err.msg
}

func NewError(msg string) *Error {
	return &Error{
		msg: msg,
	}
}
