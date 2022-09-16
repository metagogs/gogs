package session

import "github.com/metagogs/gogs/utils/slicex"

var DefaultSessionPool SessionPool

func GetSessionByID(id int64) (*Session, error) {
	return DefaultSessionPool.GetSessionByID(id)
}

// GetSessionByUID get session by user id.
// the filter is used to filter the sessions that should not receive the message.
func GetSessionByUID(uid string, filter *SessionFilter) ([]int64, []int64) {
	return DefaultSessionPool.GetSessionByUID(uid, filter)
}

// SendMessageByID send message to the session with the given id.
func SendMessageByID(sessionId int64, in interface{}) {
	if sess, err := DefaultSessionPool.GetSessionByID(sessionId); err == nil {
		_ = sess.SendMessage(in)
	}
}

func ListSessions() []*Session {
	return DefaultSessionPool.ListSessions()
}

// BroadcastMessage broadcast message to all sessions except the session with the given id.
// the filter is used to filter the sessions that should not receive the message.
func BroadcastMessage(users []string, send interface{}, filter *SessionFilter, exclude ...string) {
	for _, u := range users {
		if slicex.InSlice(u, exclude) {
			continue
		}
		if result, _ := GetSessionByUID(u, filter); len(result) > 0 {
			go SendMessageByID(result[0], send)
		}
	}
}
