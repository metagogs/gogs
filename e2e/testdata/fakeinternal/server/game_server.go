package server

import (
	"context"

	"github.com/metagogs/gogs/e2e/testdata/fakeinternal/logic/baseworld"
	"github.com/metagogs/gogs/e2e/testdata/game"

	"github.com/metagogs/gogs/e2e/testdata/fakeinternal/svc"
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

func (gogs *Server) BindUser(ctx context.Context, s *session.Session, in *game.BindUser) {
	l := baseworld.NewBindUserLogic(ctx, gogs.svcCtx, s)
	l.Handler(in)
}

func (gogs *Server) JoinWorld(ctx context.Context, s *session.Session, in *game.JoinWorld) {
	l := baseworld.NewJoinWorldLogic(ctx, gogs.svcCtx, s)
	l.Handler(in)
}

func (gogs *Server) UpdateUserInWorld(ctx context.Context, s *session.Session, in *game.UpdateUserInWorld) {
	l := baseworld.NewUpdateUserInWorldLogic(ctx, gogs.svcCtx, s)
	l.Handler(in)
}
