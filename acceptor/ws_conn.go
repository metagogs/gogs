package acceptor

import (
	"io"
	"net"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/juju/ratelimit"
)

type WSConn struct {
	state               int32
	conn                *websocket.Conn
	typ                 int
	reader              io.Reader
	info                *ConnInfo
	MaxMessagesInSecond uint
	bucket              *ratelimit.Bucket
}

func NewWSConn(conn *websocket.Conn, info *ConnInfo) *WSConn {
	wsConn := &WSConn{
		conn:  conn,
		info:  info,
		state: ConnStatusStart,
	}
	if info.BucketCapacity > 0 {
		wsConn.bucket = ratelimit.NewBucket(info.BucketFillInterval, info.BucketCapacity)
	}

	return wsConn
}

func (c *WSConn) SetCloseHandler(f func()) {
	c.conn.SetCloseHandler(func(code int, text string) error {
		f()
		return nil
	})
}

func (c *WSConn) GetInfo() *ConnInfo {
	return c.info
}

// GetNextMessage reads the next message available in the stream
func (c *WSConn) GetNextMessage() (b []byte, err error) {
	_, msgBytes, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	if c.info.BucketCapacity > 0 && c.bucket.TakeAvailable(1) < 1 {
		return nil, ErrMessageRateLimit
	}

	return msgBytes, nil
}

func (c *WSConn) Read(b []byte) (int, error) {
	if c.reader == nil {
		t, r, err := c.conn.NextReader()
		if err != nil {
			return 0, err
		}
		c.typ = t
		c.reader = r
	}
	n, err := c.reader.Read(b)
	if err != nil && err != io.EOF {
		return n, err
	} else if err == io.EOF {
		_, r, err := c.conn.NextReader()
		if err != nil {
			return 0, err
		}
		c.reader = r
	}

	return n, nil
}

func (c *WSConn) Write(b []byte) (int, error) {
	_ = c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	err := c.conn.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		return 0, err
	}

	return len(b), nil
}

func (c *WSConn) IsClosed() bool {
	return atomic.LoadInt32(&c.state) == ConnStatusClosed
}

func (c *WSConn) Close() error {
	atomic.StoreInt32(&c.state, ConnStatusClosed)
	return c.conn.Close()
}

func (c *WSConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *WSConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *WSConn) SetDeadline(t time.Time) error {
	if err := c.SetReadDeadline(t); err != nil {
		return err
	}

	return c.SetWriteDeadline(t)
}

func (c *WSConn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *WSConn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}
