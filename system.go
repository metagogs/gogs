package gogs

import (
	"context"

	"github.com/metagogs/gogs/proto"
	"github.com/metagogs/gogs/session"
)

type NetworkComponent struct {
	app *App
}

func NewNetworkComponent(app *App) *NetworkComponent {
	return &NetworkComponent{
		app: app,
	}
}

func (s *NetworkComponent) Pong(ctx context.Context, sess *session.Session, pong *proto.Pong) {
	// go s.app.LatencyServer.Pong(sess.ID(), pong.Time)
}
