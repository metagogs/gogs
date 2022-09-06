package gogs

import (
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/metagogs/gogs/acceptor"
	"github.com/metagogs/gogs/admin"
	"github.com/metagogs/gogs/agent"
	"github.com/metagogs/gogs/component"
	"github.com/metagogs/gogs/config"
	"github.com/metagogs/gogs/global"
	"github.com/metagogs/gogs/group"
	"github.com/metagogs/gogs/gslog"
	"github.com/metagogs/gogs/handler"
	"github.com/metagogs/gogs/latency"
	"github.com/metagogs/gogs/message"
	"github.com/metagogs/gogs/session"
	"github.com/metagogs/gogs/system"
	"github.com/metagogs/gogs/webserver"
	"go.uber.org/zap"
)

type App struct {
	*zap.Logger

	acceptors []acceptor.Acceptor // 网络监听
	running   bool                // 是否运行
	dieChan   chan bool           // 关闭通道

	handler       *handler.HandlerService // 消息处理
	sessionPool   session.SessionPool     // session池
	adminServer   *admin.AdminServer      // admin服务
	MessageServer *message.MessageServer  // 消息管理
	webServer     *webserver.WebServer
	LatencyServer *latency.LatencyServer // 延时服务管理
	GroupServer   *group.GroupServer     // 组管理

	helper *appHelper //一些便捷操作
	Config *config.Config
}

func NewApp(config *config.Config) *App {
	fmt.Println("============================== gogs starting ==============================")
	//设置config全局变量
	global.GoGSDebug = config.Debug
	appLog := gslog.NewLog("gogs")

	appBuidler := NewBuilder(config)

	// agent管理
	agentFactory := agent.NewAgentFactory(config, appBuidler.sessionPool, appBuidler.messageServer)

	// 消息处理
	appHandler := handler.NewHanlderService(config, agentFactory, appBuidler.messageServer)

	app := &App{
		Logger:        appLog,
		handler:       appHandler,
		dieChan:       make(chan bool),
		Config:        config,
		GroupServer:   appBuidler.groupServer,
		sessionPool:   appBuidler.sessionPool,
		adminServer:   appBuidler.adminServer,
		MessageServer: appBuidler.messageServer,
		webServer:     appBuidler.webServer,
	}

	//初始化便捷服务
	helper := newAppHelper(app)
	app.helper = helper

	system.RegisterSystemComponent(app.MessageServer, NewNetworkComponent(app))

	return app
}

func (app *App) GetSessionPool() session.SessionPool {
	return app.sessionPool
}

func (app *App) RegisterComponent(sd component.ComponentDesc, ss interface{}) {
	app.MessageServer.RegisterComponent(sd, ss)
	app.adminServer.RegisterComponent(sd, ss)
}

func (app *App) AddAcceptor(acceptor acceptor.Acceptor) {
	//todo 检测name和route是否重复
	app.acceptors = append(app.acceptors, acceptor)
	app.adminServer.AddAcceptor(acceptor)
}

func (app *App) GetAcceptors() []acceptor.Acceptor {
	return app.acceptors
}

func (app *App) SetAgentFactory(factory *agent.AgentFacotry) {
	app.handler.SetAgentFacotry(factory)
}

func (app *App) UseDefaultEncodeJSON() {
	app.MessageServer.UseDefaultEncodeJSON()
}

func (app *App) UseDefaultEncodeProto() {
	app.MessageServer.UseDefaultEncodeProto()
}

func (app *App) UseDefaultEncodeJSONWithHeader() {
	app.MessageServer.UseDefaultEncodeJSONWithHeader()
}

func (app *App) RegisterWebHandler(port int, f func(gin *gin.Engine)) {
	app.webServer.RegisterWebHandler(port, f)
}

func (app *App) Helper() *appHelper {
	return app.helper
}

func (app *App) Start() {
	app.webServer.Start()   // 监听业务http服务，包含debug测试
	app.adminServer.Start() // 启动内置admin服务，包含健康检测接口
	app.listen()            // 服务监听
	fmt.Println("============================== gogo started  ==============================")
	if app.Config.StaredCallback != nil {
		go app.Config.StaredCallback()
	}

	sg := make(chan os.Signal, 1)
	signal.Notify(sg, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	select {
	case <-app.dieChan:
		app.Info("the app will shutdown in a few seconds")
	case s := <-sg:
		app.Sugar().Info("got signal: ", s, ", shutting down...")
		close(app.dieChan)
	}

	app.sessionPool.CloseAll()
	app.Info("game server is shutting down...")
}

func (app *App) Shutdown() {
	select {
	case <-app.dieChan: // prevent closing closed channel
	default:
		close(app.dieChan)
	}
}

func (app *App) listen() {
	// start listening for connections
	for _, acc := range app.acceptors {
		a := acc
		go func() {
			for conn := range a.GetConnChan() {
				go app.handler.Handle(conn)
			}
		}()
		go func() {
			a.ListenAndServe()
		}()

		app.Sugar().Infof("listening with acceptor[%s] %s on addr %s", a.GetName(), reflect.TypeOf(a), a.GetAddr())
	}

	app.running = true
}
