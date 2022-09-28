package server

import (
	"context"

	"github.com/metagogs/gogs/e2e/smoke/internal/logic/baseworld"
	"github.com/metagogs/gogs/e2e/smoke/model"

	"github.com/metagogs/gogs/e2e/smoke/internal/svc"
	"github.com/metagogs/gogs/session"
)

type Server struct {
	svcCtx *svc.ServiceContext
}

func NewServer(svcCtx *svc.ServiceContext) *Server {
	return &Server{
		svcCtx: svcCtx,
	}
}

func (gogs *Server) BindUser(ctx context.Context, s *session.Session, in *model.BindUser) {
	l := baseworld.NewBindUserLogic(ctx, gogs.svcCtx, s)
	l.Handler(in)
}
