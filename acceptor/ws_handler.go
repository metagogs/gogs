package acceptor

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/metagogs/gogs/global"
	"github.com/metagogs/gogs/packet"
	"go.uber.org/zap"
)

type wsConnHandler struct {
	upgrader *websocket.Upgrader
	connChan chan AcceptorConn
	*zap.SugaredLogger
	info *ConnInfo
}

func (h *wsConnHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	if r.Method == "OPTIONS" {
		return
	}

	conn, err := h.upgrader.Upgrade(rw, r, nil)
	if err != nil {
		h.Warnf("Upgrade failure, URI=%s, Error=%s", r.RequestURI, err.Error())
		return
	}

	c := NewWSConn(conn, h.info)

	c.conn.SetReadLimit(packet.MaxPacketSize)

	h.connChan <- c

	if global.GoGSDebug {
		h.Info("Upgrade success",
			zap.String("name", h.info.AcceptorName),
			zap.String("group", h.info.AcceptorType),
			zap.String("host", r.Host),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("referer", r.Referer()),
			zap.String("uri", r.RequestURI))
	}
}
