package acceptor

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"sync/atomic"

	"github.com/metagogs/gogs/gslog"
	"github.com/pion/webrtc/v3"
	"go.uber.org/zap"
)

// WebRTCAcceptor accept the webrtc datachannel connection
type WebRTCAcceptor struct {
	connChan     chan AcceptorConn
	httpListener net.Listener
	udpListener  *net.UDPConn
	api          *webrtc.API
	config       *AcceptroConfig
	state        int32
	closeMutex   sync.Mutex
}

func NewWebRTCAcceptor(config *AcceptroConfig) *WebRTCAcceptor {
	if len(config.Name) == 0 {
		gslog.NewLog("webrtc_acceptor").Error("name length is 0")
		os.Exit(1)
	}
	if len(config.Groups) == 0 {
		gslog.NewLog("webrtc_acceptor").Error("groups length is 0")
		os.Exit(1)
	}
	return &WebRTCAcceptor{
		config:   config,
		connChan: make(chan AcceptorConn),
	}
}

func (w *WebRTCAcceptor) GetConfig() *AcceptroConfig {
	return w.config
}

func (w *WebRTCAcceptor) GetName() string {
	return w.config.Name
}

func (w *WebRTCAcceptor) GetConnChan() chan AcceptorConn {
	return w.connChan
}

func (w *WebRTCAcceptor) GetAddr() string {
	return fmt.Sprintf("%d[udp-webrtc] %d[http]", w.config.UdpPort, w.config.HttpPort)
}

func (w *WebRTCAcceptor) GetType() string {
	return ACCEPTOR_TYPE_WEBRTC
}

func (w *WebRTCAcceptor) Stop() {
	w.closeMutex.Lock()
	defer w.closeMutex.Unlock()

	if atomic.LoadInt32(&w.state) == StatusClosed {
		return
	}
	errUdp := w.udpListener.Close()
	if errUdp != nil {
		gslog.NewLog("webrtc_acceptor").Error("Failed to stop udp", zap.Error(errUdp))
	}
	errHttp := w.httpListener.Close()
	if errHttp != nil {
		gslog.NewLog("webrtc_acceptor").Error("Failed to stop http", zap.Error(errHttp))
	}
	if errHttp == nil && errUdp == nil {
		atomic.StoreInt32(&w.state, StatusClosed)
		gslog.NewLog("webrtc_acceptor").Info("wx_acceptr stoped")
	}
}

func (w *WebRTCAcceptor) ListenAndServe() {
	udpListener, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IP{0, 0, 0, 0},
		Port: w.config.UdpPort,
	})
	if err != nil {
		gslog.NewLog("webrtc_acceptor").Error("Failed to listen udp", zap.Error(err))
		panic(err)
	}
	w.udpListener = udpListener
	settingEngine := webrtc.SettingEngine{}
	settingEngine.DetachDataChannels()
	settingEngine.SetICEUDPMux(webrtc.NewICEUDPMux(nil, udpListener))
	w.api = webrtc.NewAPI(webrtc.WithSettingEngine(settingEngine))

	// exchange the answer and offer with the http
	httpListener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", w.config.HttpPort))
	if err != nil {
		gslog.NewLog("webrtc_acceptor").Error("Failed to listen http", zap.Error(err))
		panic(err)
	}
	w.httpListener = httpListener

	w.serve()
}

func (w *WebRTCAcceptor) serve() {
	defer w.Stop()
	_ = http.Serve(w.httpListener, &webRTCConnHandler{
		api:           w.api,
		config:        w.config,
		connChan:      w.connChan,
		SugaredLogger: gslog.NewLog("webrtc_handler").Sugar(),
		localAddr:     w.udpListener.LocalAddr(),
	})
}
