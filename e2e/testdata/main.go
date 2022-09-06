package testdata

import (
	"context"
	"time"

	"github.com/metagogs/gogs"
	"github.com/metagogs/gogs/acceptor"
	"github.com/metagogs/gogs/config"
	"github.com/metagogs/gogs/e2e/testdata/fakeinternal/server"
	"github.com/metagogs/gogs/e2e/testdata/fakeinternal/svc"
	"github.com/metagogs/gogs/e2e/testdata/game"
)

var TestApp *gogs.App

func StartServer(closeCtx context.Context, config *config.Config) {

	TestApp = gogs.NewApp(config)

	go func() {
		<-closeCtx.Done()
		TestApp.Shutdown()
	}()

	TestApp.AddAcceptor(acceptor.NewWSAcceptror(&acceptor.AcceptroConfig{
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

	TestApp.AddAcceptor(acceptor.NewWebRTCAcceptor(&acceptor.AcceptroConfig{
		HttpPort: 8889,
		UdpPort:  9001,
		Name:     "world",
		Groups: []*acceptor.AcceptorGroupConfig{
			&acceptor.AcceptorGroupConfig{
				GroupName:          "data",
				BucketFillInterval: 40 * time.Millisecond,
				BucketCapacity:     10,
			},
		},
	}))

	ctx := svc.NewServiceContext(TestApp)
	gameSrv := server.NewServer(ctx)
	game.RegisterAllComponents(TestApp, gameSrv)

	httpSrv := server.NewWebServer(ctx)
	TestApp.RegisterWebHandler(8890, httpSrv.RegisterHandler)

	defer TestApp.Shutdown()
	TestApp.Start()

	// for testing
	for _, acc := range TestApp.GetAcceptors() {
		acc.Stop()
	}

	TestApp.Info("test server is shutdown")
}
