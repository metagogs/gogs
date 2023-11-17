package acceptor

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/metagogs/gogs/gslog"
	"go.uber.org/zap"
)

type WSAcceptor struct {
	connChan   chan AcceptorConn
	listener   net.Listener
	addr       int
	config     *AcceptorConfig
	state      int32
	closeMutex sync.Mutex
}

func NewWSAcceptor(config *AcceptorConfig) *WSAcceptor {
	if len(config.Name) == 0 {
		gslog.NewLog("ws_acceptor").Error("name length is 0")
		os.Exit(1)
	}
	if len(config.Groups) == 0 {
		gslog.NewLog("ws_acceptor").Error("groups length is 0")
		os.Exit(1)
	}
	return &WSAcceptor{
		addr:     config.HttpPort,
		connChan: make(chan AcceptorConn, 50),
		config:   config,
	}
}

func (w *WSAcceptor) GetConfig() *AcceptorConfig {
	return w.config
}

func (w *WSAcceptor) GetName() string {
	return w.config.Name
}

func (w *WSAcceptor) ListenAndServe() {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", w.addr))
	if err != nil {
		gslog.NewLog("ws_acceptor").Error("Failed to listen", zap.Error(err))
		panic(err)
	}
	w.listener = listener

	w.serve(&upgrader)
}

func (w *WSAcceptor) Stop() {
	w.closeMutex.Lock()
	defer w.closeMutex.Unlock()

	if atomic.LoadInt32(&w.state) == StatusClosed {
		return
	}
	err := w.listener.Close()
	if err != nil {
		gslog.NewLog("ws_acceptor").Error("Failed to stop", zap.Error(err))
		return
	}
	atomic.StoreInt32(&w.state, StatusClosed)
	gslog.NewLog("ws_acceptor").Info("ws_acceptor stoped")
}

func (w *WSAcceptor) GetConnChan() chan AcceptorConn {
	return w.connChan
}

func (w *WSAcceptor) GetAddr() string {
	return fmt.Sprintf("%d[http]", w.config.HttpPort)
}

func (w *WSAcceptor) GetType() string {
	return ACCEPTOR_TYPE_WS
}

func adaptWSHandler(handler *wsConnHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})
}

func (w *WSAcceptor) serve(upgrader *websocket.Upgrader) {
	defer w.Stop()

	mux := http.NewServeMux()
	for _, group := range w.config.Groups {

		handler := &wsConnHandler{
			upgrader:      upgrader,
			connChan:      w.connChan,
			SugaredLogger: gslog.NewLog("ws_handler").Sugar(),
			info: &ConnInfo{
				AcceptorType:       w.GetType(),
				AcceptorName:       w.GetName(),
				AcceptorGroup:      group.GroupName,
				BucketFillInterval: group.BucketFillInterval,
				BucketCapacity:     group.BucketCapacity,
			},
		}

		adaptedHandler := adaptWSHandler(handler)

		for _, middleware := range group.MiddlewareFunc {
			adaptedHandler = middleware(adaptedHandler)
		}

		mux.Handle("/"+group.GroupName, adaptedHandler)
	}
	_ = http.Serve(w.listener, mux)
}
