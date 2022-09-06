package dispatch

import (
	"context"
	"os"
	"reflect"
	"strconv"
	"sync"

	"github.com/metagogs/gogs/component"
	"github.com/metagogs/gogs/gslog"
	"github.com/metagogs/gogs/packet"
	"github.com/metagogs/gogs/session"
)

// manage method and filed
type DispatchServer struct {
	mu        sync.Mutex
	methods   map[uint32]*serverMethod
	objAction map[string]uint32
}

func NewDispatchServer() *DispatchServer {
	return &DispatchServer{
		methods:   make(map[uint32]*serverMethod),
		objAction: make(map[string]uint32),
	}
}

func (a *DispatchServer) RegisterComponent(sd component.ComponentDesc, ss interface{}) {
	gslog.NewLog("server").Sugar().Infof("RegisterComponent: %s", sd.ComonentName)
	if ss != nil {
		ht := reflect.TypeOf(sd.ComponentType).Elem()
		st := reflect.TypeOf(ss)
		if !st.Implements(ht) {
			gslog.NewLog("server").Sugar().Warnf("gogs: Server.RegisterComponent found the handler of type %v that does not satisfy %v", st, ht)
			os.Exit(1)
		}
	}
	a.register(sd, ss)
}

func (a *DispatchServer) register(sd component.ComponentDesc, ss interface{}) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for i := range sd.Methods {
		d := &sd.Methods[i]
		d.MethodName = d.FieldType.Name()
		serverMethod := &serverMethod{
			srv:           ss,
			componentDesc: &sd,
			methodDesc:    d,
			fieldHandler:  d.FiledHanler,
		}
		if d.FieldType.Kind() == reflect.Ptr {
			serverMethod.fieldType = d.FieldType.Elem()
		} else {
			serverMethod.fieldType = d.FieldType
		}

		if _, ok := a.objAction[d.FieldType.Name()]; ok {
			gslog.NewLog("dispatch").Sugar().Warnf("gogs: Server.RegisterComponent found duplicate field type %s", d.FieldType.Name())
			os.Exit(1)
		}

		gslog.NewLog("dispatch").Sugar().Infof("RegisterFieldAction: %s %d", d.FieldType.Name(), d.MethodIndex)
		a.objAction[d.FieldType.Name()] = d.MethodIndex
		a.methods[d.MethodIndex] = serverMethod

		if d.Handler != nil {
			gslog.NewLog("dispatch").Sugar().Infof("RegisterMethod: %s with key %s", d.MethodName, strconv.FormatUint(uint64(d.MethodIndex), 10))
		}
	}
}

func (app *DispatchServer) GetMethod(actionKey uint32) (*serverMethod, bool) {
	serverMethod, ok := app.methods[actionKey]
	return serverMethod, ok
}

func (app *DispatchServer) GetAction(in interface{}) (uint32, error) {
	t := reflect.TypeOf(in)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	action, ok := app.objAction[t.Name()]
	if !ok {
		return 0, ErrObjectActionNotFound
	}

	return action, nil
}

func (app *DispatchServer) GetActionByName(fieldName string) (uint32, error) {
	actionNum, ok := app.objAction[fieldName]
	if !ok {
		return 0, ErrObjectActionNotFound
	}

	return actionNum, nil
}

func (app *DispatchServer) GetObjType(actionKey uint32) (reflect.Type, bool) {
	method, ok := app.GetMethod(actionKey)
	if !ok {
		return nil, false
	}
	if method.fieldType == nil {
		return nil, false
	}

	return method.fieldType, true
}

func (app *DispatchServer) GetObj(actionKey uint32) (interface{}, bool) {
	method, ok := app.GetMethod(actionKey)
	if !ok {
		return nil, false
	}

	return method.fieldHandler(), true
}

func (app *DispatchServer) Call(ctx context.Context, sess *session.Session, packet *packet.Packet) error {
	method, ok := app.GetMethod(packet.GetActionKey())
	if !ok {
		return ErrMethodNotFound
	}
	if method.methodDesc.Handler == nil {
		return ErrMethodNotFound
	}
	method.Call(ctx, sess, packet.Bean)

	return nil
}
