package model

import (
	"context"
	"reflect"

	"github.com/metagogs/gogs"
	"github.com/metagogs/gogs/component"
	"github.com/metagogs/gogs/packet"
	"github.com/metagogs/gogs/session"
)

func RegisterAllComponents(s *gogs.App, srv Component) {
	registerBaseWorldComponent(s, srv)

}

func registerBaseWorldComponent(s *gogs.App, srv Component) {
	s.RegisterComponent(_BaseWorldComponentDesc, srv)
}

type Component interface {
	BindUser(ctx context.Context, s *session.Session, in *BindUser)
}

func _BaseWorldComponent_BindUser_Handler(srv interface{}, ctx context.Context, sess *session.Session, in interface{}) {
	srv.(Component).BindUser(ctx, sess, in.(*BindUser))
}

var _BaseWorldComponentDesc = component.ComponentDesc{
	ComponentName:  "BaseWorldComponent",
	ComponentIndex: 1, // equal to module index
	ComponentType:  (*Component)(nil),
	Methods: []component.ComponentMethodDesc{
		{
			MethodIndex: packet.CreateAction(packet.ServicePacket, 1, 1), // 0x810001 8454145
			FieldType:   reflect.TypeOf(BindUser{}),
			Handler:     _BaseWorldComponent_BindUser_Handler,
			FieldHandler: func() interface{} {
				return new(BindUser)
			},
		},
		{
			MethodIndex: packet.CreateAction(packet.ServicePacket, 1, 2), // 0x810002 8454146
			FieldType:   reflect.TypeOf(BindSuccess{}),
			Handler:     nil,
			FieldHandler: func() interface{} {
				return new(BindSuccess)
			},
		},
	},
}
