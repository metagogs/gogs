package player

import (
	"sync"
)

type PlayerManager struct {
	players sync.Map
}

func NewPlayerManager() *PlayerManager {
	return &PlayerManager{}
}

func (p *PlayerManager) CreateUser(uid, name string) {
	player := newPlayer(uid, name, p)
	p.players.Store(uid, player)
}

func (p *PlayerManager) GetPlayer(uid string) (*Player, bool) {
	v, ok := p.players.Load(uid)
	if !ok {
		return nil, false
	}

	return v.(*Player), true
}

func (p *PlayerManager) DeletePlayer(uid string) {
	p.players.Delete(uid)
}
