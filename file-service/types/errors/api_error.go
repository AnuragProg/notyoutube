package errors

import (
	"fmt"
	"runtime/debug"
)

type APIError struct {
	Inner      error
	Message    string
	StatusCode int
	StackTrace string
}

func NewAPIError(statusCode int, messagef string, msgArgs ...interface{}) APIError {
	return APIError{
		Inner: nil,
		Message: fmt.Sprintf(messagef, msgArgs...),
		StatusCode: statusCode,
		StackTrace: string(debug.Stack()),
	}
}

func IntoAPIError(err error, statusCode int, messagef string, msgArgs ...interface{}) APIError {
	return APIError{
		Inner:   err,
		Message: fmt.Sprintf(messagef, msgArgs...),
		StatusCode: statusCode,
		StackTrace: string(debug.Stack()),
	}
}

func (err APIError) Error() string {
	return err.Message
}
