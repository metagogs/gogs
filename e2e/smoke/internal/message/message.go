package message

import (
	"github.com/metagogs/gogs/e2e/smoke/model"
	"github.com/metagogs/gogs/session"
)

func SendBindUser(s *session.Session, in *model.BindUser) error {
	return s.SendMessage(in, "BindUser")
}

func SendBindSuccess(s *session.Session, in *model.BindSuccess) error {
	return s.SendMessage(in, "BindSuccess")
}
