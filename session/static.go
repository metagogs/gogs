package session

var DefaultSessionPool SessionPool

func GetSessionByID(id int64) (*Session, error) {
	return DefaultSessionPool.GetSessionByID(id)
}

func GetSessionByUID(uid string, filter *SessionFilter) ([]int64, []int64) {
	return DefaultSessionPool.GetSessionByUID(uid, filter)
}

func ListSessions() []*Session {
	return DefaultSessionPool.ListSessions()
}
