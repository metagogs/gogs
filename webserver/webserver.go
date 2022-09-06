package webserver

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/metagogs/gogs/config"
)

// WebServer
// Provide the web http server with different port
// Use the gin framework
type WebServer struct {
	config    *config.Config
	httpServe map[int]*gin.Engine // http服务
}

func NewWebServer(config *config.Config) *WebServer {
	return &WebServer{
		httpServe: make(map[int]*gin.Engine),
		config:    config,
	}
}

// RegisterWebHandler register web handler, use the custom port and callback the gin engine for user
func (app *WebServer) RegisterWebHandler(port int, f func(gin *gin.Engine)) {
	if _, exist := app.httpServe[port]; exist {
		panic("port is used")
	}
	app.httpServe[port] = app.craeteServer()
	f(app.httpServe[port])
}

// createServer create the gin engine
func (app *WebServer) craeteServer() *gin.Engine {
	return gin.New()
}

func (app *WebServer) Start() {
	for port, g := range app.httpServe {
		go func(e *gin.Engine, p int) {
			_ = e.Run(fmt.Sprintf("0.0.0.0:%d", p))
		}(g, port)
	}
	if app.config.Debug {
		// in the debug mode show the pprof handle
		app.debugServer()
	}
}

// debugServer the pprof handle
func (app *WebServer) debugServer() {
	pp := http.NewServeMux()
	pp.HandleFunc("/debug/pprof/", pprof.Index)
	pp.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	pp.HandleFunc("/debug/pprof/profile", pprof.Profile)
	pp.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	pp.HandleFunc("/debug/pprof/trace", pprof.Trace)
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.GopprofAddr),
		Handler:           pp,
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		_ = server.ListenAndServe()
	}()
}
