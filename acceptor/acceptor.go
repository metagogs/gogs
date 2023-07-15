package acceptor

import (
	"net"
	"time"
)

type ConnInfo struct {
	AcceptorType       string        // the type of the acceptor
	AcceptorName       string        // the name of the acceptor
	AcceptorGroup      string        // the group of the acceptor
	Ordered            bool          // whether the message is ordered, only for webrtc
	BucketFillInterval time.Duration // the interval of the bucket fill
	BucketCapacity     int64         // the capacity of the bucket
}

// AcceptorConn iface
type AcceptorConn interface {
	GetNextMessage() (b []byte, err error) // get the next message
	GetInfo() *ConnInfo                    // get the conn info
	SetCloseHandler(func())                // set the close callback
	IsClosed() bool                        // is conn closed
	net.Conn                               // Conn
}

// Acceptor type interface
type Acceptor interface {
	GetConfig() *AcceptorConfig     // get the config
	ListenAndServe()                // listen and serve
	Stop()                          // stop the acceptor
	GetConnChan() chan AcceptorConn // get the conn channel
	GetAddr() string                // get the addr
	GetName() string                // get the name
	GetType() string                // get the type
}
