package custom_error

import (
	"fmt"
)

type Error struct {
	HTTPCode int
	Code     int
	Message  string
}

func (c Error) Error() string {
	return fmt.Sprintf("%d - %s", c.Code, c.Message)
}

func NewError(httpCode int, code int, message string) error {
	return Error{
		HTTPCode: httpCode,
		Code:     code,
		Message:  message,
	}
}
