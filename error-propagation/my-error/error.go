package error

import (
	"fmt"
	"runtime/debug"
)

// MyError is used to represent detailed errors
type MyError struct {
	Inner      error
	Message    string
	StackTrace string
	Misc       map[string]interface{}
}

// WrapError constructs a MyError struct
func WrapError(err error, messageF string, msgArgs ...interface{}) MyError {
	return MyError{
		Inner:      err,
		Message:    fmt.Sprintf(messageF, msgArgs...),
		StackTrace: string(debug.Stack()),
		Misc:       make(map[string]interface{}),
	}
}

// method on MyError struct that returns its error message
func (err MyError) Error() string {
	return err.Message
}
