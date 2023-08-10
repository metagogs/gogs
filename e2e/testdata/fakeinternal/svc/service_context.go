package svc

import (
	"fmt"
	"os"

	"github.com/bwmarrin/snowflake"
	"github.com/metagogs/gogs"
	"github.com/metagogs/gogs/e2e/testdata/fakeinternal/player"
	"github.com/metagogs/gogs/group"
	"github.com/metagogs/gogs/utils/snow"
)

type ServiceContext struct {
	*gogs.App
	SF            *snowflake.Node
	PlayerManager *player.PlayerManager
	World         *group.MemoryGroup
}

func NewServiceContext(app *gogs.App) *ServiceContext {
	sf, err := snow.NewSnowNode()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	pl := player.NewPlayerManager()
	//test
	pl.CreateUser("123", "neosu")
	pl.CreateUser("456", "szp")

	world := group.NewMemoryGroup("world", sf.Generate().Int64())

	return &ServiceContext{
		App:           app,
		SF:            sf,
		PlayerManager: pl,
		World:         world,
	}
}
