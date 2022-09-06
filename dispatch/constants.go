package dispatch

import "errors"

var (
	ErrObjectActionNotFound = errors.New("object action not found")
	ErrMethodNotFound       = errors.New("method not found")
)
