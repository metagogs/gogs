package gogs

import "github.com/metagogs/gogs/session"

type appHelper struct {
	*App
}

func newAppHelper(app *App) *appHelper {
	return &appHelper{
		App: app,
	}
}

func (h *appHelper) GetSessionByID(id int64) (*session.Session, error) {
	return h.sessionPool.GetSessionByID(id)
}

func (h *appHelper) GetSessionByUID(uid string, filter *session.SessionFilter) ([]int64, []int64) {
	return h.sessionPool.GetSessionByUID(uid, filter)
}

func (h *appHelper) SendMessageByID(sessionId int64, in interface{}) {
	if sess, err := h.GetSessionByID(sessionId); err == nil {
		_ = sess.SendMessage(in)
	}
}
