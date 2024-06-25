package gogs

import (
	"flag"
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
	"github.com/metagogs/gogs/deployment"
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

	Config *config.Config
}

func NewApp(config *config.Config) *App {
	fmt.Println("============================== gogs starting ==============================")
	//设置config全局变量
	global.GoGSDebug = config.Debug
	appLog := gslog.NewLog("gogs")

	appBuilder := NewBuilder(config)

	// agent管理
	agentFactory := agent.NewAgentFactory(config, appBuilder.sessionPool, appBuilder.messageServer)

	// 消息处理
	appHandler := handler.NewHandlerService(config, agentFactory, appBuilder.messageServer)

	app := &App{
		Logger:        appLog,
		handler:       appHandler,
		dieChan:       make(chan bool),
		Config:        config,
		GroupServer:   appBuilder.groupServer,
		sessionPool:   appBuilder.sessionPool,
		adminServer:   appBuilder.adminServer,
		MessageServer: appBuilder.messageServer,
		webServer:     appBuilder.webServer,
	}

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

func (app *App) SetAgentFactory(factory *agent.AgentFactory) {
	app.handler.SetAgentFactory(factory)
}

func (app *App) UseDefaultEncodeJSON() {
	app.MessageServer.UseDefaultEncodeJSON()
}

func (app *App) UseDefaultEncodeProto() {
	app.MessageServer.UseDefaultEncodeProto()
}

// UseDefaultEncodePureJSON 发送消息时不包含头部,使用JSON,方便测试使用,纯JSON
// 编码有两种，一种是带头部的，一种是不带头部的
// PureJSON是不带头部的，类型是写在action里面的
func (app *App) UseDefaultEncodePureJSON() {
	app.MessageServer.UseDefaultEncodePureJSON()
}

func (app *App) RegisterWebHandler(port int, f func(gin *gin.Engine)) {
	app.webServer.RegisterWebHandler(port, f)
}

func (app *App) Start() {
	deploymentFlag := flag.Bool("deployment", false, "deployment mode")
	deploymentScv := flag.Bool("svc", false, "use the k8s svc")
	deploymentName := flag.String("name", "", "deployment name")
	deploymentSpace := flag.String("namespace", "", "deployment namespace")
	flag.Parse()
	if *deploymentFlag {
		deploymentHelper := deployment.NewDeploymentHelper(app.Config, *deploymentScv, *deploymentName, *deploymentSpace)
		for _, acc := range app.acceptors {
			deploymentHelper.AddAcceptor(acc.GetConfig())
		}
		// generate deployment file
		_ = deploymentHelper.Generate()
		return
	}

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
