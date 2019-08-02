package uranus

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	// ErrContextNil is thrown if the context passed is nil to any function that require to use context.
	ErrContextNil = errors.New("Context is Nil")

	// ErrNotFound is thrown if the item requested does not exist in the system
	ErrNotFound = errors.New("Your requested item does not exists")

	// ErrNotModified is thrown to the client when the cached copy of a particular file is up to date with the server.
	ErrNotModified = errors.New("")
)

// ConstraintError represents a custom error for a contstraint things.
type ConstraintError string

func (e ConstraintError) Error() string {
	return string(e)
}

// ConstraintErrorf constructs ConstraintError with formatted message.
func ConstraintErrorf(format string, a ...interface{}) ConstraintError {
	return ConstraintError(fmt.Sprintf(format, a...))
}

// ErrorFromResponseStatusCode generates error based on the status code from *http.Response.
// For example, it will generate fetlar.ErrNotFound when given status code of 404.
func ErrorFromResponseStatusCode(code int, message string) (err error) {
	switch code {
	case http.StatusNotFound:
		err = ErrNotFound
	case http.StatusBadRequest:
		err = ConstraintErrorf(message)
	case http.StatusNotModified:
		err = ErrNotModified
	default:
		err = fmt.Errorf(message)
	}

	return
}
