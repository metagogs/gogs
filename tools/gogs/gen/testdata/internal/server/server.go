package server

import (
	"context"

	"github.com/metagogs/test/internal/logic/baseworld"
	"github.com/metagogs/test/model"

	"github.com/metagogs/gogs/session"
	"github.com/metagogs/test/internal/svc"
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
