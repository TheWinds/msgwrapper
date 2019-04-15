package msgwrapper

type StatusError interface {
	Code() int32
	Message() string
	Detail() string
	From(error) StatusError
	error
}

type Error struct {
	code    int32
	message string
	detail  string
}

func NewError(code int32, message string) *Error {
	return &Error{code: code, message: message}
}

func NewErrorWithDetail(code int32, message string, detail string) *Error {
	return &Error{code: code, message: message, detail: detail}
}

func (e *Error) Code() int32 {
	return e.code
}

func (e *Error) From(err error) StatusError {
	if err == nil {
		return nil
	}
	if isStatusError(err) {
		return err.(StatusError)
	}
	return NewError(500, err.Error())
}

func (e *Error) Message() string {
	return e.message
}

func (e *Error) Detail() string {
	return e.detail
}

func (e *Error) Error() string {
	return e.message
}

func isStatusError(err error) bool {
	_, ok := err.(StatusError)
	return ok
}
