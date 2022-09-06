package packet

import (
	"sync"
)

var (
	_pool = NewPool()

	GetPacket = _pool.GetPacket
	PutPacket = _pool.PutPacket
)

type Pool struct {
	packet *sync.Pool
}

func NewPool() *Pool {
	p := &Pool{
		packet: &sync.Pool{
			New: func() interface{} {
				return &Packet{
					HeaderByte: make([]byte, 8),
				}
			},
		},
	}

	return p
}

func (p Pool) GetPacket() *Packet {
	return p.packet.Get().(*Packet)
}

func (p Pool) PutPacket(bean *Packet) {
	p.packet.Put(bean)
}
