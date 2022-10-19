package session

import (
	"github.com/metagogs/gogs/acceptor"
	"github.com/metagogs/gogs/networkentity"
	"github.com/metagogs/gogs/packet"
	"go.uber.org/zap"
)

type Session struct {
	uid              string                      // session/user id
	id               int64                       // agent id
	pool             *sessionPoolImpl            // session pool
	sessionLog       *zap.Logger                 // session log
	agent            networkentity.NetworkEntity // agent
	data             SessionData                 // session data
	OnCloseCallbacks []func(id int64)            // close callback
}

func (sess *Session) Close() {
	sess.log().Info("session close")
	sess.pool.DeleteSession(sess.id)
	if len(sess.uid) != 0 {
		sess.pool.deleteSessionByUID(sess.uid, sess.id)
	}
	_ = sess.agent.Stop()
}

func (sess *Session) log() *zap.Logger {
	return sess.sessionLog.With(zap.Int64("agent_id", sess.id))
}

func (sess *Session) ID() int64 {
	return sess.agent.GetId()
}

func (sess *Session) UID() string {
	return sess.uid
}

func (sess *Session) SetUID(uid string) {
	sess.uid = uid
	sess.pool.addSessionByUID(uid, sess)
}

func (sess *Session) IsLogin() bool {
	return len(sess.uid) > 0
}

func (sess *Session) SendMessage(in interface{}, name ...string) error {
	if sess.pool.config.SendMessageLog {
		sess.log().Info("send message")
	}
	return sess.agent.Send(in, name...)
}

func (sess *Session) SendData(data []byte) {
	sess.agent.SendData(data)
}

func (sess *Session) SendPacket(data *packet.Packet) {
	sess.agent.SendPacket(data)
}

func (sess *Session) GetLastTimeOnline() int64 {
	return sess.agent.GetLastTimeOnline()
}

func (sess *Session) GetData() SessionData {
	return sess.data
}

func (sess *Session) OnClose(c func(id int64)) error {
	sess.OnCloseCallbacks = append(sess.OnCloseCallbacks, c)
	return nil
}

func (sess *Session) GetOnCloseCallbacks() []func(id int64) {
	return sess.OnCloseCallbacks
}

func (sess *Session) SetOnCloseCallbacks(callbacks []func(id int64)) {
	sess.OnCloseCallbacks = callbacks
}

func (sess *Session) SetOnCloseCallback(c func(id int64)) {
	sess.OnCloseCallbacks = []func(id int64){c}
}

func (sess *Session) GetLatency() int64 {
	return sess.agent.GetLatency()
}

func (sess *Session) GetConnInfo() *acceptor.ConnInfo {
	return sess.agent.GetConnInfo()
}
