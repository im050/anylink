package errs

import "fmt"

const (
	// E 错误
	E = 50000
)

// Error HTTP错误
type Error struct {
	Message string
	Code    int
}

func (e Error) Error() string {
	return e.Message
}

// New 返回一个新的错误
func New(message string, codes ...int) *Error {
	code := E
	if len(codes) > 0 {
		code = codes[0]
	}
	return &Error{
		Message: message,
		Code:    code,
	}
}

func (h *Error) Output(message string) *Error {
	h.Message = message
	return h
}

func Errorf(tpl string, args ...string) *Error {
	return New(fmt.Sprintf(tpl, args))
}
