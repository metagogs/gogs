package session

import "sync"

type SessionMemory struct {
	sync.RWMutex
	data map[string]interface{}
}

func NewSessionMemory() *SessionMemory {
	return &SessionMemory{
		data: make(map[string]interface{}),
	}
}

// GetData gets the data
func (sess *SessionMemory) GetData() map[string]interface{} {
	sess.RLock()
	defer sess.RUnlock()

	return sess.data
}

// Set associates value with the key in session storage
func (sess *SessionMemory) Set(key string, value interface{}) {
	sess.Lock()
	defer sess.Unlock()

	sess.data[key] = value
}

func (sess *SessionMemory) Get(key string) (interface{}, bool) {
	sess.RLock()
	defer sess.RUnlock()

	value, ok := sess.data[key]
	return value, ok
}

func (sess *SessionMemory) Delete(key string) {
	sess.RLock()
	defer sess.RUnlock()

	delete(sess.data, key)
}

func (sess *SessionMemory) GetString(key, def string) string {
	sess.RLock()
	defer sess.RUnlock()

	value, ok := sess.data[key]
	if ok {
		return value.(string)
	}

	return def
}
