package baseworld

import (
	"context"

	"github.com/metagogs/gogs/e2e/smoke/internal/svc"
	"github.com/metagogs/gogs/e2e/smoke/model"
	"github.com/metagogs/gogs/gslog"
	"github.com/metagogs/gogs/session"
	"go.uber.org/zap"
)

type BindUserLogic struct {
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	session *session.Session
	*zap.Logger
}

func NewBindUserLogic(ctx context.Context, svcCtx *svc.ServiceContext, sess *session.Session) *BindUserLogic {
	return &BindUserLogic{
		ctx:     ctx,
		svcCtx:  svcCtx,
		session: sess,
		Logger:  gslog.NewLog("bind_user_logic"),
	}
}

func (l *BindUserLogic) Handler(in *model.BindUser) {

}
