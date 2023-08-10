package group

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/bwmarrin/snowflake"
	"github.com/metagogs/gogs/utils/snow"
)

type GroupServer struct {
	rooms sync.Map
	sf    *snowflake.Node
}

func NewGroupServer() *GroupServer {
	sf, err := snow.NewSnowNode()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return &GroupServer{
		sf: sf,
	}
}

func (gs *GroupServer) CreateMemoryGroup(name string) Group {
	if group, ok := gs.rooms.Load(name); ok {
		return group.(*MemoryGroup)
	}

	group := NewMemoryGroup(name, gs.sf.Generate().Int64())
	gs.rooms.Store(name, group)

	return group
}

func (gs *GroupServer) GetGroup(name string) (Group, bool) {
	if group, ok := gs.rooms.Load(name); ok {
		return group.(Group), true
	}

	return nil, false
}

func (gs *GroupServer) DeleteGroup(name string) {
	gs.rooms.Delete(name)
}

func (gs *GroupServer) DeleteUserByName(uid string) {
	gs.rooms.Range(func(key, value interface{}) bool {
		if group, ok := value.(Group); ok {
			_ = group.RemoveUser(context.TODO(), uid)
		}
		return true
	})
}

func (gs *GroupServer) ListGroup() []Group {
	list := []Group{}
	gs.rooms.Range(func(key, value any) bool {
		list = append(list, value.(Group))
		return true
	})

	return list
}
