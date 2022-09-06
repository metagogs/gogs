package {{.LogicPackage}}

import (
	"context"

	"{{.BasePackage}}/{{.Package}}"
	"{{.BasePackage}}/internal/svc"
	"github.com/metagogs/gogs/gslog"
	"github.com/metagogs/gogs/session"
	"go.uber.org/zap"
)

type {{.Name}}Logic struct {
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	session *session.Session
	*zap.Logger
}

func New{{.Name}}Logic(ctx context.Context, svcCtx *svc.ServiceContext, sess *session.Session) *{{.Name}}Logic {
	return &{{.Name}}Logic{
		ctx:     ctx,
		svcCtx:  svcCtx,
		session: sess,
		Logger:  gslog.NewLog("{{.SnakeName}}_logic"),
	}
}

func (l *{{.Name}}Logic) Handler(in *{{.Package}}.{{.Name}}) {

}
