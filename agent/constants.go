package agent

import "errors"

var (
	ErrCloseClosedSession = errors.New("close closed session")
)

const (
	_ int32 = iota
	// StatusStart status
	StatusStart
	// StatusHandshake status
	StatusHandshake
	// StatusWorking status
	StatusWorking
	// StatusClosed status
	StatusClosed
)
