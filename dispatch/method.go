package dispatch

import (
	"context"
	"reflect"

	"github.com/metagogs/gogs/component"
	"github.com/metagogs/gogs/session"
)

type serverMethod struct {
	srv           interface{}
	componentDesc *component.ComponentDesc
	methodDesc    *component.ComponentMethodDesc
	fieldType     reflect.Type
	fieldHandler  func() interface{}
}

func (s *serverMethod) GetSrv() interface{} {
	return s.srv
}

func (s *serverMethod) GetMethodDesc() *component.ComponentMethodDesc {
	return s.methodDesc
}

func (s *serverMethod) String() string {
	return s.methodDesc.MethodName
}

func (s *serverMethod) NewType() interface{} {
	if s.fieldType == nil {
		return nil
	}

	return reflect.New(s.fieldType).Interface()
}

func (s *serverMethod) Call(ctx context.Context, sess *session.Session, in interface{}) {
	s.methodDesc.Handler(s.srv, ctx, sess, in)
}
