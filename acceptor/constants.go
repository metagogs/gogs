package acceptor

import "errors"

var (
	ACCEPTOR_TYPE_WS     = "websockets"
	ACCEPTOR_TYPE_WEBRTC = "webrtc"

	ErrMessageRateLimit = errors.New("message rate limit")
)

const (
	_ int32 = iota
	ConnStatusStart
	ConnStatusClosed
)

const (
	_ int32 = iota
	StatusWorking
	StatusClosed
)
