package admin

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/metagogs/gogs/acceptor"
	"github.com/metagogs/gogs/component"
	"github.com/metagogs/gogs/config"
	"github.com/metagogs/gogs/gslog"
	"github.com/metagogs/gogs/session"
	"go.uber.org/zap"
)

// AdminServer get the system running info
type AdminServer struct {
	*gin.Engine
	*zap.SugaredLogger
	addr        string
	debug       bool
	config      *config.Config
	sessionPool session.SessionPool
	decodeType  map[string]uint8
	encodeType  string
	status      systemStatus
}

func NewAdminServer(config *config.Config,
	sessionPool session.SessionPool,
	decodeType map[string]uint8,
	encodeType string) *AdminServer {

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	server := &AdminServer{
		Engine:        r,
		addr:          fmt.Sprintf("0.0.0.0:%d", config.AdminPort),
		debug:         config.Debug,
		config:        config,
		SugaredLogger: gslog.NewLog("admin").Sugar(),
		sessionPool:   sessionPool,
		decodeType:    decodeType,
		encodeType:    encodeType,
		status: systemStatus{
			Acceptors: make(map[string]string),
		},
	}
	server.Init()

	return server
}

func (a *AdminServer) Init() {
	a.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// for k8s health probe
	a.Any("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	if a.debug {
		// in the debug mode, show the more system info
		a.GET("/admin", a.systemStatus)
	}
}

func (a *AdminServer) Start() {
	a.Infof("admin server start on %s", a.addr)
	go func() {
		_ = a.Run(a.addr)
	}()
}

func (a *AdminServer) systemStatus(g *gin.Context) {
	a.Info("request system status")

	a.status.DebugMode = a.config.Debug
	a.status.DecodeType = a.decodeType
	a.status.EncodeType = a.encodeType
	a.status.SessionCount = a.sessionPool.GetSessionTotalCount()  // total session connected include disconnect in the system
	a.status.OnlineSessionCount = a.sessionPool.GetSessionCount() // current connected session count
	a.status.NumGoroutine = runtime.NumGoroutine()
	a.status.NumCPU = runtime.NumCPU()
	a.status.NumCgoCall = runtime.NumCgoCall()

	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	a.status.Memory = systemMemory{
		Alloc:        m.Alloc,
		TotalAlloc:   m.TotalAlloc,
		Sys:          m.Sys,
		Lookups:      m.Lookups,
		Mallocs:      m.Mallocs,
		Frees:        m.Frees,
		HeapAlloc:    m.HeapAlloc,
		HeapSys:      m.HeapSys,
		HeapIdle:     m.HeapIdle,
		HeapInuse:    m.HeapInuse,
		HeapReleased: m.HeapReleased,
		HeapObjects:  m.HeapObjects,
		StackInuse:   m.StackInuse,
		StackSys:     m.StackSys,
	}

	a.status.Env = os.Environ()

	g.JSON(200, a.status)
}

func (a *AdminServer) RegisterComponent(sd component.ComponentDesc, ss interface{}) {
	a.status.Components = append(a.status.Components, sd)
}

func (a *AdminServer) AddAcceptor(acceptor acceptor.Acceptor) {
	a.status.Acceptors[reflect.TypeOf(acceptor).Elem().Name()] = acceptor.GetAddr()
}
