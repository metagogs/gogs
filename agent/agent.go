package agent

import (
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/metagogs/gogs/acceptor"
	"github.com/metagogs/gogs/message"
	"github.com/metagogs/gogs/packet"
	"github.com/metagogs/gogs/proto"
	"github.com/metagogs/gogs/session"
	"github.com/metagogs/gogs/utils/bytebuffer"
	"go.uber.org/zap"
)

// Agent corresponds to a user and is used for storing raw Conn information
type Agent struct {
	AgentID          int64                       // user id
	conn             acceptor.AcceptorConn       // the Conn
	agentLog         *zap.Logger                 // logger
	chSend           chan *bytebuffer.ByteBuffer // send message channel
	chSendByte       chan []byte                 // send message byte channel
	chStopWrite      chan struct{}               // close message send
	chStopHeartbeat  chan struct{}               // close heartbeat
	chDie            chan struct{}               // wait for close
	closeMutex       sync.Mutex                  // close mutes
	sess             *session.Session            // the session which bind to this agent
	messageServer    *message.MessageServer      // message server can decode and encode the message and call the handler
	heartbeatTimeout time.Duration               // heartbeat timeout
	heartbeatLog     bool                        // heartbeat log
	lastAt           int64                       // last heartbeat time
	state            int32                       // agent state
	pingPool         *sync.Pool                  // ping sync pool
}

func (a *Agent) GetId() int64 {
	return a.AgentID
}

func (a *Agent) GetSession() *session.Session {
	return a.sess
}

func (a *Agent) Start() {
	a.log().Info("agent start")
	// start write message channel
	go a.write()
	// start heartbeat
	go a.heartbeat()
	a.conn.SetCloseHandler(func() {
		_ = a.Stop()
	})
	<-a.chDie
}

// Stop close the agent
func (a *Agent) Stop() error {
	a.closeMutex.Lock()
	defer a.closeMutex.Unlock()

	// check if already closed
	if a.GetStatus() == StatusClosed {
		return ErrCloseClosedSession
	}
	a.SetStatus(StatusClosed)

	select {
	case <-a.chDie:
	default:
		close(a.chStopWrite)
		close(a.chStopHeartbeat)
		close(a.chDie)
		a.onSessionClosed()
	}

	return a.conn.Close()
}

func (a *Agent) IsClosed() bool {
	return a.GetStatus() == StatusClosed
}

func (a *Agent) onSessionClosed() {
	defer func() {
		if err := recover(); err != nil {
			a.log().Error("onSessionClosed error", zap.Any("recover", err))
		}
	}()

	for _, fn := range a.sess.GetOnCloseCallbacks() {
		fn(a.sess.ID())
	}

}

// write write the message from the channel to the Conn
func (a *Agent) write() {
	defer func() {
		_ = a.Stop()
	}()

	for {
		select {
		case msg := <-a.chSend:
			if _, err := a.conn.Write(msg.B); err != nil {
				bytebuffer.Put(msg)
				return
			}
			bytebuffer.Put(msg)
		case msg := <-a.chSendByte:
			if _, err := a.conn.Write(msg); err != nil {
				return
			}
		case <-a.chStopWrite:
			return
		}
	}
}

func (a *Agent) heartbeat() {
	if a.heartbeatTimeout == 0 {
		// if heartbeat timeout is 0, then close the heartbeat
		return
	}
	ticker := time.NewTicker(a.heartbeatTimeout)

	defer func() {
		ticker.Stop()
		_ = a.Stop()
	}()

	for {
		select {
		case <-ticker.C:
			if a.heartbeatLog {
				deadline := time.Now().Add(-2 * a.heartbeatTimeout).Unix()
				if atomic.LoadInt64(&a.lastAt) < deadline {
					a.log().Warn("Session heartbeat timeout",
						zap.Int64("LastTime", atomic.LoadInt64(&a.lastAt)),
						zap.Int64("Deadline", deadline))
				}
			}
			a.sendHeartbeat()
		case <-a.chDie:
			return
		case <-a.chStopHeartbeat:
			return
		}
	}
}

// sendHeartbeat message
func (a *Agent) sendHeartbeat() {
	heartTime := time.Now().UnixMilli()
	heartTimeStr := strconv.FormatInt(heartTime, 10)

	ping, _ := pingPool.Get().(*proto.Ping)
	ping.Time = heartTimeStr
	defer pingPool.Put(ping)

	if a.heartbeatLog {
		a.log().Info("send heartbeat")
	}
	_ = a.Send(ping)
}

// Send Message
func (a *Agent) Send(in interface{}, name ...string) error {
	// encode the message with the message server
	data, err := a.messageServer.EncodeMessage(in, name...)
	if err != nil {
		a.log().Error("encode error", zap.Error(err))
		return err
	}

	select {
	case a.chSend <- data.ToData():
	case <-a.chDie:
	}

	return nil
}

func (a *Agent) SendPacket(data *packet.Packet) {
	select {
	case a.chSend <- data.ToData():
	case <-a.chDie:
	}
}

func (a *Agent) SendData(data []byte) {
	select {
	case a.chSendByte <- data:
	case <-a.chDie:
	}
}

func (a *Agent) SetLastAt() {
	atomic.StoreInt64(&a.lastAt, time.Now().Unix())
}

func (a *Agent) GetLastTimeOnline() int64 {
	return atomic.LoadInt64(&a.lastAt)
}

func (a *Agent) GetStatus() int32 {
	return atomic.LoadInt32(&a.state)
}

func (a *Agent) SetStatus(state int32) {
	atomic.StoreInt32(&a.state, state)
}

func (a *Agent) log() *zap.Logger {
	return a.agentLog
}

func (a *Agent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *Agent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *Agent) GetLatency() int64 {
	//todo get latency
	return 0
}

func (a *Agent) GetConnInfo() *acceptor.ConnInfo {
	return a.conn.GetInfo()
}
