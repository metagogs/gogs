package message

import (
	"github.com/metagogs/gogs/session"
	"github.com/metagogs/test/model"
)

func SendBindUser(s *session.Session, in *model.BindUser) error {
	return s.SendMessage(in, "BindUser")
}

func SendBindSuccess(s *session.Session, in *model.BindSuccess) error {
	return s.SendMessage(in, "BindSuccess")
}
