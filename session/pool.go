package session

import (
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/metagogs/gogs/config"
	"github.com/metagogs/gogs/gslog"
	"github.com/metagogs/gogs/networkentity"
	"github.com/metagogs/gogs/utils/slicex"
	"go.uber.org/zap"
)

type SessionPool interface {
	CreateSession(agent networkentity.NetworkEntity) *Session  // create session
	DeleteSession(id int64)                                    // delete session by id
	GetSessionCount() int64                                    // get current session count
	GetSessionTotalCount() int64                               // get total session count
	GetSessionByID(int64) (*Session, error)                    // get session by id
	GetSessionByUID(string, *SessionFilter) ([]int64, []int64) // get session by uid with filter
	ListSessions() []*Session                                  // list all session
	CloseAll()                                                 // close all session
}

// SessionFilter can filter the session from session list
// one user maybe has many session, so we can filter the session
type SessionFilter struct {
	ConnType  string
	ConnName  string
	ConnGroup string
}

type sessionList struct {
	mutex sync.Mutex
	data  map[int64]*sessionWrap
	list  []int64
}

func (s *sessionList) Add(w *sessionWrap) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data[w.SessionID] = w
	s.list = append(s.list, w.SessionID)
}

func (s *sessionList) Delete(id int64) int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.data, id)
	s.list = slicex.RemoveSliceItem(s.list, id)
	return len(s.data)
}

// GetList get session list by filter, return the filter result and all the sessions as backup
// if we can not get the result, we can check user's all the session list
// then we can send message use another session
func (s *sessionList) GetList(filter *SessionFilter) ([]int64, []int64) {
	if filter == nil {
		return s.list, s.list
	}
	result := []int64{}
	for _, v := range s.data {
		if filter != nil {
			if len(filter.ConnType) > 0 && v.ConnType != filter.ConnType {
				continue
			}

			if len(filter.ConnName) > 0 && v.ConnName != filter.ConnName {
				continue
			}

			if len(filter.ConnGroup) > 0 && v.ConnGroup != filter.ConnGroup {
				continue
			}
		}

		result = append(result, v.SessionID)
	}

	return result, s.list
}

type sessionWrap struct {
	SessionID int64
	ConnType  string
	ConnName  string
	ConnGroup string
}

type sessionPoolImpl struct {
	config       *config.Config
	currentCount int64
	count        int64
	sessionsByID sync.Map
	sessionByUID sync.Map
	sessionLog   *zap.Logger
}

func NewSessionPool(config *config.Config) *sessionPoolImpl {
	return &sessionPoolImpl{
		config:     config,
		sessionLog: gslog.NewLog("session_pool"),
	}
}

func (pool *sessionPoolImpl) CreateSession(agent networkentity.NetworkEntity) *Session {
	s := &Session{
		id:         agent.GetId(),
		uid:        strconv.FormatInt(agent.GetId(), 10),
		agent:      agent,
		pool:       pool,
		sessionLog: gslog.NewLog("session"),
		data:       NewSessionMemeory(),
	}

	pool.sessionsByID.Store(s.id, s)
	atomic.AddInt64(&pool.currentCount, 1)
	atomic.AddInt64(&pool.count, 1)

	pool.sessionLog.Info("session created",
		zap.Int64("current_count", pool.GetSessionCount()),
		zap.Int64("id", s.id))

	return s
}

func (pool *sessionPoolImpl) addSessionByUID(uid string, sess *Session) {
	warp := &sessionWrap{
		SessionID: sess.id,
		ConnType:  sess.GetConnInfo().AcceptorType,
		ConnName:  sess.GetConnInfo().AcceptorName,
		ConnGroup: sess.GetConnInfo().AcceptorGroup,
	}

	if v, ok := pool.sessionByUID.Load(uid); ok {
		v.(*sessionList).Add(warp)
	} else {
		list := &sessionList{
			data: make(map[int64]*sessionWrap),
		}
		list.Add(warp)
		pool.sessionByUID.Store(uid, list)
	}
}

func (pool *sessionPoolImpl) deleteSessionByUID(uid string, id int64) {
	if v, ok := pool.sessionByUID.Load(uid); ok {
		if v.(*sessionList).Delete(id) == 0 {
			pool.sessionByUID.Delete(uid)
		}
	}
}

func (pool *sessionPoolImpl) DeleteSession(id int64) {
	atomic.AddInt64(&pool.currentCount, -1)
	pool.sessionsByID.Delete(id)

	pool.sessionLog.Info("session deleted",
		zap.Int64("current_count", pool.GetSessionCount()),
		zap.Int64("id", id))
}

func (pool *sessionPoolImpl) GetSessionCount() int64 {
	return pool.currentCount
}

func (pool *sessionPoolImpl) GetSessionTotalCount() int64 {
	return pool.count
}

func (pool *sessionPoolImpl) GetSessionByID(id int64) (*Session, error) {
	if v, ok := pool.sessionsByID.Load(id); ok {
		return v.(*Session), nil
	}

	return nil, ErrSessionNotFound
}

func (pool *sessionPoolImpl) GetSessionByUID(uid string, filter *SessionFilter) ([]int64, []int64) {
	if v, ok := pool.sessionByUID.Load(uid); ok {
		return v.(*sessionList).GetList(filter)
	}

	return nil, nil
}

func (pool *sessionPoolImpl) CloseAll() {
	pool.sessionsByID.Range(func(key, value any) bool {
		if s, ok := value.(*Session); ok {
			s.Close()
		}
		return true
	})
	pool.sessionLog.Info("Close all session")
}

func (pool *sessionPoolImpl) ListSessions() []*Session {
	result := []*Session{}
	pool.sessionsByID.Range(func(key, value any) bool {
		if s, ok := value.(*Session); ok {
			result = append(result, s)
		}
		return true
	})
	return result
}
