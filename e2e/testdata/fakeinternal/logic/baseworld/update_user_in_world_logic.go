package baseworld

import (
	"context"

	"github.com/metagogs/gogs/e2e/testdata/fakeinternal/message"
	"github.com/metagogs/gogs/e2e/testdata/fakeinternal/svc"
	"github.com/metagogs/gogs/e2e/testdata/game"
	"github.com/metagogs/gogs/session"
)

type UpdateUserInWorldLogic struct {
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	session *session.Session
}

func NewUpdateUserInWorldLogic(ctx context.Context, svcCtx *svc.ServiceContext, sess *session.Session) *UpdateUserInWorldLogic {
	return &UpdateUserInWorldLogic{
		ctx:     ctx,
		svcCtx:  svcCtx,
		session: sess,
	}
}

func (l *UpdateUserInWorldLogic) Handler(in *game.UpdateUserInWorld) {
	if !l.session.IsLogin() {
		return
	}

	player, exist := l.svcCtx.PlayerManagaer.GetPlayer(l.session.UID())
	if !exist {
		return
	}
	//make sure uid is right
	in.Uid = player.UID

	uids := l.svcCtx.World.GetUsers(l.ctx)
	for _, u := range uids {
		if u != l.session.UID() {
			go l.notify(u, in)
		}
	}
}

func (l *UpdateUserInWorldLogic) notify(uid string, send *game.UpdateUserInWorld) {

	if result, _ := l.svcCtx.Helper().GetSessionByUID(uid, nil); len(result) > 0 {
		if sess, err := l.svcCtx.Helper().GetSessionByID(result[0]); err == nil {
			_ = message.SendUpdateUserInWorld(sess, send)
		}
	}
}
