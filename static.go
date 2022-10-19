package gogs

import (
	"github.com/metagogs/gogs/message"
	"github.com/metagogs/gogs/packet"
	"github.com/metagogs/gogs/session"
	"github.com/metagogs/gogs/utils/slicex"
)

var DefaultSessionPool session.SessionPool
var DefaultMessageServer *message.MessageServer

func ListSessions() []*session.Session {
	return DefaultSessionPool.ListSessions()
}

func GetSessionByID(id int64) (*session.Session, error) {
	return DefaultSessionPool.GetSessionByID(id)
}

// GetSessionByUID get session by user id.
// the filter is used to filter the sessions that should not receive the message.
func GetSessionByUID(uid string, filter *session.SessionFilter) ([]int64, []int64) {
	return DefaultSessionPool.GetSessionByUID(uid, filter)
}

// SendMessageByID send message to the session with the given id.
func SendMessageByID(sessionId int64, in interface{}) {
	if sess, err := DefaultSessionPool.GetSessionByID(sessionId); err == nil {
		_ = sess.SendMessage(in)
	}
}

func SendDataByID(sessionId int64, in []byte) {
	if sess, err := DefaultSessionPool.GetSessionByID(sessionId); err == nil {
		sess.SendData(in)
	}
}

func SendPacketByID(sessionId int64, in *packet.Packet) {
	if sess, err := DefaultSessionPool.GetSessionByID(sessionId); err == nil {
		sess.SendPacket(in)
	}
}

// BroadcastMessage broadcast message to all sessions except the session with the given id.
// the filter is used to filter the sessions that should not receive the message.
// send pakcet
func BroadcastMessage(users []string, send interface{}, filter *session.SessionFilter, exclude ...string) error {
	packet, err := EncodeMessage(send)
	if err != nil {
		return err
	}
	for _, u := range users {
		if slicex.InSlice(u, exclude) {
			continue
		}
		if result, _ := GetSessionByUID(u, filter); len(result) > 0 {
			SendPacketByID(result[0], packet)
		}
	}

	return nil
}

// send bytes
func BroadcastData(users []string, data []byte, filter *session.SessionFilter, exclude ...string) {
	for _, u := range users {
		if slicex.InSlice(u, exclude) {
			continue
		}
		if result, _ := GetSessionByUID(u, filter); len(result) > 0 {
			SendDataByID(result[0], data)
		}
	}
}

func EncodeMessage(in interface{}, name ...string) (*packet.Packet, error) {
	packet, err := DefaultMessageServer.EncodeMessage(in, name...)
	if err != nil {
		return nil, err
	}
	return packet, nil
}
