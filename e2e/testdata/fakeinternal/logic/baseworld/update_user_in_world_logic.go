package baseworld

import (
	"context"

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
	session.BroadcastMessage(uids, in, nil, l.session.UID())
}
