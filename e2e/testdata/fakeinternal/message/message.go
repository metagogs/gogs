package message

import (
	"github.com/metagogs/gogs/e2e/testdata/game"
	"github.com/metagogs/gogs/session"
)

func SendBindUser(s *session.Session, in *game.BindUser) error {
	return s.SendMessage(in, "BindUser")
}

func SendBindSuccess(s *session.Session, in *game.BindSuccess) error {
	return s.SendMessage(in, "BindSuccess")
}

func SendJoinWorld(s *session.Session, in *game.JoinWorld) error {
	return s.SendMessage(in, "JoinWorld")
}

func SendJoinWorldSuccess(s *session.Session, in *game.JoinWorldSuccess) error {
	return s.SendMessage(in, "JoinWorldSuccess")
}

func SendJoinWorldNotify(s *session.Session, in *game.JoinWorldNotify) error {
	return s.SendMessage(in, "JoinWorldNotify")
}

func SendUpdateUserInWorld(s *session.Session, in *game.UpdateUserInWorld) error {
	return s.SendMessage(in, "UpdateUserInWorld")
}
