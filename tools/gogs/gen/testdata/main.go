package main

import (
	"github.com/metagogs/gogs"
	"github.com/metagogs/gogs/acceptor"
	"github.com/metagogs/gogs/config"
	"github.com/metagogs/test/internal/server"
	"github.com/metagogs/test/internal/svc"
	"github.com/metagogs/test/model"
)

func main() {

	config := config.NewConfig()

	app := gogs.NewApp(config)
	app.AddAcceptor(acceptor.NewWSAcceptror(&acceptor.AcceptroConfig{
		HttpPort: 8888,
		Name:     "base",
		Groups: []*acceptor.AcceptorGroupConfig{
			&acceptor.AcceptorGroupConfig{
				GroupName: "base",
			},
		},
	}))

	app.AddAcceptor(acceptor.NewWebRTCAcceptor(&acceptor.AcceptroConfig{
		HttpPort: 8889,
		UdpPort:  9000,
		Name:     "world",
		Groups: []*acceptor.AcceptorGroupConfig{
			&acceptor.AcceptorGroupConfig{
				GroupName: "data",
			},
		},
	}))

	ctx := svc.NewServiceContext(app)
	srv := server.NewServer(ctx)

	model.RegisterAllComponents(app, srv)

	defer app.Shutdown()
	app.Start()
}
