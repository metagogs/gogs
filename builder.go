package gogs

import (
	"github.com/metagogs/gogs/admin"
	"github.com/metagogs/gogs/codec"
	"github.com/metagogs/gogs/config"
	"github.com/metagogs/gogs/dispatch"
	"github.com/metagogs/gogs/group"
	"github.com/metagogs/gogs/message"
	"github.com/metagogs/gogs/session"
	"github.com/metagogs/gogs/webserver"
)

type builder struct {
	sessionPool   session.SessionPool    // session池管理
	groupServer   *group.GroupServer     // 组管理
	adminServer   *admin.AdminServer     // 服务管理
	messageServer *message.MessageServer // 消息管理，包装编码和分发
	webServer     *webserver.WebServer   // web服务管理
}

func NewBuilder(config *config.Config) *builder {
	webserver := newWebServer(config)
	sessionPool := newSessionPool(config)
	groupServer := newGroupServer(config)
	dispatchServer := newDispatchServer(config)
	codecHelper := newCodecHelper(config, dispatchServer)
	adminServer := newAdminServer(config, sessionPool, codecHelper)
	messageServer := newMessageServer(config, codecHelper, dispatchServer)

	return &builder{
		sessionPool:   sessionPool,
		groupServer:   groupServer,
		adminServer:   adminServer,
		messageServer: messageServer,
		webServer:     webserver,
	}
}

func newSessionPool(config *config.Config) session.SessionPool {
	// session池管理
	appSessionPool := session.NewSessionPool(config)
	session.DefaultSessionPool = appSessionPool
	return appSessionPool
}

func newDispatchServer(config *config.Config) *dispatch.DispatchServer {
	// 消息和方法分发管理
	return dispatch.NewDispatchServer()
}

func newGroupServer(config *config.Config) *group.GroupServer {
	return group.NewGroupServer()
}

func newCodecHelper(config *config.Config, dispatchServer *dispatch.DispatchServer) *codec.CodecHelper {
	return codec.NewCodecHelper(config, dispatchServer)
}

func newAdminServer(config *config.Config, sessionPool session.SessionPool, codecHelper *codec.CodecHelper) *admin.AdminServer {
	d, e := codecHelper.GetTypes()
	return admin.NewAdminServer(config, sessionPool, d, e)
}

func newMessageServer(config *config.Config, codecHelper *codec.CodecHelper, dispatchServer *dispatch.DispatchServer) *message.MessageServer {
	return message.NewMessageServer(codecHelper, dispatchServer)
}

func newWebServer(config *config.Config) *webserver.WebServer {
	return webserver.NewWebServer(config)
}
