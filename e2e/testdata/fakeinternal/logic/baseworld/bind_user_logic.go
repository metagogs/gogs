package baseworld

import (
	"context"
	"strings"

	"github.com/metagogs/gogs/e2e/testdata/fakeinternal/message"
	"github.com/metagogs/gogs/e2e/testdata/fakeinternal/svc"
	"github.com/metagogs/gogs/e2e/testdata/game"
	"github.com/metagogs/gogs/session"
)

var BindUserHandler = make(chan *game.BindUser)

type BindUserLogic struct {
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	session *session.Session
}

func NewBindUserLogic(ctx context.Context, svcCtx *svc.ServiceContext, sess *session.Session) *BindUserLogic {
	return &BindUserLogic{
		ctx:     ctx,
		svcCtx:  svcCtx,
		session: sess,
	}
}

func (l *BindUserLogic) Handler(in *game.BindUser) {
	BindUserHandler <- in
	if strings.HasPrefix(in.Uid, "test_") {
		l.svcCtx.PlayerManagaer.CreateUser(in.Uid, in.Uid)
	}
	player, ok := l.svcCtx.PlayerManagaer.GetPlayer(in.Uid)
	if ok {
		l.session.SetUID(in.Uid)
		l.session.GetData().Set("name", player.Name)
		player.AddSession(l.session.GetConnInfo().AcceptorName, l.session.GetConnInfo().AcceptorGroup, l.session.ID())
		l.session.OnClose(func(id int64) {
			player.DeleteSession(l.session.ID())
		})

		_ = message.SendBindSuccess(l.session, &game.BindSuccess{})

	} else {
	}
}
