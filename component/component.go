package component

import (
	"context"
	"reflect"

	"github.com/metagogs/gogs/session"
)

// methodHandler is the handler of component method
//
//	func _BaseWorldComponent_BindUser_Handler(srv interface{}, ctx context.Context, sess *session.Session, in interface{}) {
//		srv.(Component).BindUser(ctx, sess, in.(*BindUser))
//	}
type methodHandler func(srv interface{}, ctx context.Context, sess *session.Session, in interface{})

// return the field instance
// it's better to use the filed handler to create the filed instance
//
//	func() interface{} {
//			return new(BindUser)
//	}
type fieldHandler func() interface{}

type ComponentDesc struct {
	ComponentName  string
	ComponentIndex uint8       // equal to module index
	ComponentType  interface{} `json:"-"`
	Methods        []ComponentMethodDesc
}

type ComponentMethodDesc struct {
	MethodIndex  uint32        // equal to action index
	MethodName   string        // equal to action name
	FieldType    reflect.Type  `json:"-"` // method field type to create field by reflect, but use the filed handler is better
	Handler      methodHandler `json:"-"` // method handler
	FieldHandler fieldHandler  `json:"-"` // method field handler, the handler function will be called to create field
}

type Component interface{}
