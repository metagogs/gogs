package baseworld

import (
	"context"
	"sync"

	"github.com/metagogs/gogs"
	"github.com/metagogs/gogs/e2e/testdata/fakeinternal/message"
	"github.com/metagogs/gogs/e2e/testdata/fakeinternal/svc"
	"github.com/metagogs/gogs/e2e/testdata/game"
	"github.com/metagogs/gogs/session"
)

var JoinWorldHandler = make(chan *game.JoinWorld)

type JoinWorldLogic struct {
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	session *session.Session
}

func NewJoinWorldLogic(ctx context.Context, svcCtx *svc.ServiceContext, sess *session.Session) *JoinWorldLogic {
	return &JoinWorldLogic{
		ctx:     ctx,
		svcCtx:  svcCtx,
		session: sess,
	}
}

var beanPool = sync.Pool{
	New: func() interface{} {
		return &game.JoinWorldNotify{}
	},
}

func (l *JoinWorldLogic) Handler(in *game.JoinWorld) {
	JoinWorldHandler <- in
	if !l.session.IsLogin() {
		return
	}

	player, exist := l.svcCtx.PlayerManager.GetPlayer(l.session.UID())
	if !exist {
		return
	}

	if err := l.svcCtx.World.AddUser(l.ctx, player.UID); err != nil {
	}
	player.OnExist(func() {
		_ = l.svcCtx.World.RemoveUser(l.ctx, player.UID)
	})

	worldUids := l.svcCtx.World.GetUsers(l.ctx)

	_ = message.SendJoinWorldSuccess(l.session, &game.JoinWorldSuccess{
		Uids: worldUids,
	})

	sendMsg, _ := beanPool.Get().(*game.JoinWorldNotify)
	sendMsg.Uid = player.UID
	sendMsg.Name = player.Name
	defer beanPool.Put(sendMsg)

	uids := l.svcCtx.World.GetUsers(l.ctx)
	gogs.BroadcastMessage(uids, sendMsg, nil, l.session.UID())
}
