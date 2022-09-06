package acceptor

import (
	"io"
	"math"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/juju/ratelimit"
	"github.com/metagogs/gogs/gslog"
	"github.com/pion/datachannel"
	"github.com/pion/webrtc/v3"
	"go.uber.org/zap"
)

var rlBufPool = sync.Pool{New: func() interface{} {
	return make([]byte, math.MaxUint16) // message size limit for Chromium
}}

type WebRTCConn struct {
	state        int32 //状态
	info         *ConnInfo
	connection   *webrtc.PeerConnection
	dataChannel  *webrtc.DataChannel
	rw           datachannel.ReadWriteCloser
	remoteAddr   net.Addr
	localAddr    net.Addr
	reader       io.Reader
	closeHandler func()
	bucket       *ratelimit.Bucket
}

func NewWebRTCConn(conn *webrtc.PeerConnection,
	dataChannel *webrtc.DataChannel,
	rw datachannel.ReadWriteCloser,
	info *ConnInfo) *WebRTCConn {

	webRTRConn := &WebRTCConn{
		info:        info,
		connection:  conn,
		dataChannel: dataChannel,
		rw:          rw,
		state:       ConnStatusStart,
	}
	if info.BucketCapacity > 0 {
		webRTRConn.bucket = ratelimit.NewBucket(info.BucketFillInterval, info.BucketCapacity)
	}
	webRTRConn.start()

	return webRTRConn
}

func (c *WebRTCConn) start() {
	c.connection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		if connectionState == webrtc.ICEConnectionStateDisconnected {
			c.closeConnection()
		} else if connectionState == webrtc.ICEConnectionStateFailed {
			c.closeConnection()
		} else if connectionState == webrtc.ICEConnectionStateClosed {
			c.closeConnection()
		}
	})
	c.connection.OnConnectionStateChange(func(pcs webrtc.PeerConnectionState) {
		if pcs == webrtc.PeerConnectionStateClosed {
			c.closeConnection()
		}
	})
	c.dataChannel.OnClose(func() {
		c.closeConnection()
	})
}

func (c *WebRTCConn) closeConnection() {
	defer func() {
		if err := recover(); err != nil {
			gslog.NewLog("webrtc_conn").Error("close error", zap.Any("recover", err))
		}
	}()
	_ = c.connection.Close()
	atomic.StoreInt32(&c.state, ConnStatusClosed)
	c.Close()
	if c.closeHandler != nil {
		c.closeHandler()
	}

}

func (c *WebRTCConn) SetCloseHandler(f func()) {
	c.closeHandler = f
}

func (c *WebRTCConn) GetInfo() *ConnInfo {
	return c.info
}

func (c *WebRTCConn) GetNextMessage() (b []byte, err error) {
	buffer, _ := rlBufPool.Get().([]byte)
	defer rlBufPool.Put(buffer) //nolint
	n, _, err := c.rw.ReadDataChannel(buffer)
	if err != nil {
		return nil, err
	}

	m := make([]byte, n)
	copy(m, buffer[:n])

	if c.info.BucketCapacity > 0 && c.bucket.TakeAvailable(1) < 1 {
		return nil, ErrMessageRateLimit
	}

	return m, nil
}

func (c *WebRTCConn) Read(b []byte) (int, error) {
	n, _, err := c.rw.ReadDataChannel(b)
	return n, err
}

func (c *WebRTCConn) Write(b []byte) (int, error) {
	return c.rw.Write(b)
}

func (c *WebRTCConn) IsClosed() bool {
	return atomic.LoadInt32(&c.state) == ConnStatusClosed
}

func (c *WebRTCConn) Close() error {
	defer func() {
		if err := recover(); err != nil {
			gslog.NewLog("webrtc_conn").Error("close error", zap.Any("recover", err))
		}
	}()

	if err := c.dataChannel.Close(); err != nil {
		return err
	}
	atomic.StoreInt32(&c.state, ConnStatusClosed)
	if err := c.rw.Close(); err != nil {
		return err
	}
	return nil
}

func (c *WebRTCConn) LocalAddr() net.Addr {
	return c.localAddr
}

func (c *WebRTCConn) RemoteAddr() net.Addr {
	return c.remoteAddr
}

func (c *WebRTCConn) SetDeadline(t time.Time) error {
	return nil
}

func (c *WebRTCConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *WebRTCConn) SetWriteDeadline(t time.Time) error {
	return nil
}
