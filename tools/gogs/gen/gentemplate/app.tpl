package main

import (
	"time"

	"{{.BasePackage}}/{{.Package}}"
	"{{.BasePackage}}/internal/server"
	"{{.BasePackage}}/internal/svc"
	"github.com/metagogs/gogs"
	"github.com/metagogs/gogs/acceptor"
	"github.com/metagogs/gogs/config"
)

func main() {

	config := config.NewConfig()

	app := gogs.NewApp(config)
	app.AddAcceptor(acceptor.NewWSAcceptror(&acceptor.AcceptroConfig{
		HttpPort: 8888,
		Name:     "base",
		Groups: []*acceptor.AcceptorGroupConfig{
			&acceptor.AcceptorGroupConfig{
				GroupName:          "base",
				BucketFillInterval: 40 * time.Millisecond,
				BucketCapacity:     10,
			},
		},
	}))

	app.AddAcceptor(acceptor.NewWebRTCAcceptor(&acceptor.AcceptroConfig{
		HttpPort: 8889,
		UdpPort:  9000,
		Name:     "world",
		Groups: []*acceptor.AcceptorGroupConfig{
			&acceptor.AcceptorGroupConfig{
				GroupName:          "data",
				BucketFillInterval: 40 * time.Millisecond,
				BucketCapacity:     10,
			},
		},
	}))

	ctx := svc.NewServiceContext(app)
	srv := server.NewServer(ctx)

	{{.Package}}.RegisterAllComponents(app, srv)

	defer app.Shutdown()
	app.Start()
}
