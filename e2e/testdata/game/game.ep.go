package game

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

	JoinWorld(ctx context.Context, s *session.Session, in *JoinWorld)

	UpdateUserInWorld(ctx context.Context, s *session.Session, in *UpdateUserInWorld)
}

func _BaseWorldComponent_BindUser_Handler(srv interface{}, ctx context.Context, sess *session.Session, in interface{}) {
	srv.(Component).BindUser(ctx, sess, in.(*BindUser))
}

func _BaseWorldComponent_JoinWorld_Handler(srv interface{}, ctx context.Context, sess *session.Session, in interface{}) {
	srv.(Component).JoinWorld(ctx, sess, in.(*JoinWorld))
}

func _BaseWorldComponent_UpdateUserInWorld_Handler(srv interface{}, ctx context.Context, sess *session.Session, in interface{}) {
	srv.(Component).UpdateUserInWorld(ctx, sess, in.(*UpdateUserInWorld))
}

var _BaseWorldComponentDesc = component.ComponentDesc{
	ComonentName:   "BaseWorldComponent",
	ComponentIndex: 1, // equeal to module index
	ComponentType:  (*Component)(nil),
	Methods: []component.ComponentMethodDesc{
		{
			MethodIndex: packet.CreateAction(packet.ServicePacket, 1, 1),
			FieldType:   reflect.TypeOf(BindUser{}),
			Handler:     _BaseWorldComponent_BindUser_Handler,
			FiledHanler: func() interface{} {
				return new(BindUser)
			},
		},
		{
			MethodIndex: packet.CreateAction(packet.ServicePacket, 1, 2),
			FieldType:   reflect.TypeOf(JoinWorld{}),
			Handler:     _BaseWorldComponent_JoinWorld_Handler,
			FiledHanler: func() interface{} {
				return new(JoinWorld)
			},
		},
		{
			MethodIndex: packet.CreateAction(packet.ServicePacket, 1, 3),
			FieldType:   reflect.TypeOf(JoinWorldNotify{}),
			Handler:     nil,
			FiledHanler: func() interface{} {
				return new(JoinWorldNotify)
			},
		},
		{
			MethodIndex: packet.CreateAction(packet.ServicePacket, 1, 4),
			FieldType:   reflect.TypeOf(UpdateUserInWorld{}),
			Handler:     _BaseWorldComponent_UpdateUserInWorld_Handler,
			FiledHanler: func() interface{} {
				return new(UpdateUserInWorld)
			},
		},
		{
			MethodIndex: packet.CreateAction(packet.ServicePacket, 1, 5),
			FieldType:   reflect.TypeOf(BindSuccess{}),
			Handler:     nil,
			FiledHanler: func() interface{} {
				return new(BindSuccess)
			},
		},
		{
			MethodIndex: packet.CreateAction(packet.ServicePacket, 1, 6),
			FieldType:   reflect.TypeOf(JoinWorldSuccess{}),
			Handler:     nil,
			FiledHanler: func() interface{} {
				return new(JoinWorldSuccess)
			},
		},
	},
}
