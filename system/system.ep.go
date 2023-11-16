package system

import (
	"context"
	"reflect"

	"github.com/metagogs/gogs/component"
	"github.com/metagogs/gogs/message"
	"github.com/metagogs/gogs/packet"
	"github.com/metagogs/gogs/proto"
	"github.com/metagogs/gogs/session"
)

func RegisterSystemComponent(s *message.MessageServer, srv NetworkComponent) {
	s.RegisterComponent(_NetworkComponentDesc, srv)
}

type NetworkComponent interface {
	Pong(ctx context.Context, sess *session.Session, pong *proto.Pong)
}

func _NetworkComponent_Pong_Handler(srv interface{}, ctx context.Context, sess *session.Session, in interface{}) {
	srv.(NetworkComponent).Pong(ctx, sess, in.(*proto.Pong))
}

var _NetworkComponentDesc = component.ComponentDesc{
	ComponentName:  "NetworkComponent",
	ComponentIndex: 1, // equal to module index
	ComponentType:  (*NetworkComponent)(nil),
	Methods: []component.ComponentMethodDesc{
		{
			MethodIndex: packet.CreateAction(packet.SystemPacket, 1, 1),
			FieldType:   reflect.TypeOf(proto.Ping{}),
			Handler:     nil,
			FieldHandler: func() interface{} {
				return new(proto.Ping)
			},
		},
		{
			MethodIndex: packet.CreateAction(packet.SystemPacket, 1, 2),
			FieldType:   reflect.TypeOf(proto.Pong{}),
			Handler:     _NetworkComponent_Pong_Handler,
			FieldHandler: func() interface{} {
				return new(proto.Pong)
			},
		},
	},
}
