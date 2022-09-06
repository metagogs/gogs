package session

import "sync"

type SessionMemeory struct {
	sync.RWMutex
	data map[string]interface{}
}

func NewSessionMemeory() *SessionMemeory {
	return &SessionMemeory{
		data: make(map[string]interface{}),
	}
}

// GetData gets the data
func (sess *SessionMemeory) GetData() map[string]interface{} {
	sess.RLock()
	defer sess.RUnlock()

	return sess.data
}

// Set associates value with the key in session storage
func (sess *SessionMemeory) Set(key string, value interface{}) {
	sess.Lock()
	defer sess.Unlock()

	sess.data[key] = value
}

func (sess *SessionMemeory) Get(key string) (interface{}, bool) {
	sess.RLock()
	defer sess.RUnlock()

	value, ok := sess.data[key]
	return value, ok
}

func (sess *SessionMemeory) Delete(key string) {
	sess.RLock()
	defer sess.RUnlock()

	delete(sess.data, key)
}

func (sess *SessionMemeory) GetString(key, def string) string {
	sess.RLock()
	defer sess.RUnlock()

	value, ok := sess.data[key]
	if ok {
		return value.(string)
	}

	return def
}
