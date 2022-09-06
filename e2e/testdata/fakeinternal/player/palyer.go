package player

import (
	"sync"
)

type Player struct {
	UID              string
	Name             string
	Sessions         sync.Map
	SessionID        sync.Map
	playerManager    *PlayerManager
	OnCloseCallbacks []func()
}

func newPlayer(uid string, name string, playerManager *PlayerManager) *Player {
	return &Player{
		UID:           uid,
		Name:          name,
		playerManager: playerManager,
	}
}

func (player *Player) AddSession(name, group string, id int64) {
	player.Sessions.Store(name+group, id)
	player.SessionID.Store(id, name+group)
}

func (player *Player) GetSession(name, group string) (int64, bool) {
	v, ok := player.Sessions.Load(name + group)
	if !ok {
		return 0, false
	}

	return v.(int64), true
}

func (player *Player) DeleteSession(id int64) {
	if v, ok := player.SessionID.Load(id); ok {
		player.Sessions.Delete(v.(string))
		player.SessionID.Delete(id)
	}
	count := 0
	player.SessionID.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	if count == 0 {
		player.Level()
	}
}

func (player *Player) Level() {
	player.playerManager.DeletePlayer(player.UID)
	for _, fn := range player.OnCloseCallbacks {
		fn()
	}
}

func (player *Player) OnExist(f func()) {
	player.OnCloseCallbacks = append(player.OnCloseCallbacks, f)
}
