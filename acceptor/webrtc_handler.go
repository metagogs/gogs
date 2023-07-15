package acceptor

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/pion/webrtc/v3"
	"go.uber.org/zap"
)

type webRTCConnHandler struct {
	api       *webrtc.API
	config    *AcceptorConfig
	connChan  chan AcceptorConn
	localAddr net.Addr
	*zap.SugaredLogger
}

func (h *webRTCConnHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	if r.Method == "OPTIONS" {
		return
	}

	// one webrtc client is current server, so we do not need the turn server
	// because the server ip is public and fixed
	peerConnection, err := h.api.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		h.Errorw("failed to create peer connection", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = rw.Write([]byte("failed to create peer connection"))
		return
	}

	var offer webrtc.SessionDescription
	if err := json.NewDecoder(r.Body).Decode(&offer); err != nil {
		h.Errorw("failed to decode offer", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = rw.Write([]byte("failed to decode offer"))
		return
	}
	if len(offer.SDP) == 0 {
		h.Errorw("offer is empty")
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = rw.Write([]byte("offer is empty"))
		return
	}

	if err := peerConnection.SetRemoteDescription(offer); err != nil {
		h.Errorw("failed to set remote description", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = rw.Write([]byte("failed to set remote description"))
		return
	}

	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		h.Errorw("failed to create answer", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = rw.Write([]byte("failed to create answer"))
		return
	}

	if err := peerConnection.SetLocalDescription(answer); err != nil {
		h.Errorw("failed to set local description", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = rw.Write([]byte("failed to set local description"))
		return
	}

	<-gatherComplete

	response, err := json.Marshal(*peerConnection.LocalDescription())
	if err != nil {
		h.Errorw("failed to marshal local description", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = rw.Write([]byte("failed to marshal local description"))
		return
	}

	go h.WaitForConnection(peerConnection, r)

	rw.Header().Set("Content-Type", "application/json")
	if _, err := rw.Write(response); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		h.Errorw("failed to write response", zap.Error(err))
		return
	}

}

// WaitForConnection the webrtc connection can have many datachannel, is different with websockets
// so we can create datachannels for each group in one connection
// but please remember, we must have least one datachannel in one connection
func (handler *webRTCConnHandler) WaitForConnection(peer *webrtc.PeerConnection, r *http.Request) {
	for _, group := range handler.config.Groups {
		dataChannel, err := peer.CreateDataChannel(group.GroupName, &webrtc.DataChannelInit{
			Ordered: &group.Ordered,
		})
		if err != nil {
			handler.Error("failed to create data channel", zap.String("name", group.GroupName), zap.Bool("order", group.Ordered), zap.Error(err))
		}

		dataChannel.OnOpen(func() {
			handler.Info("data channel opened ", zap.String("name", group.GroupName), zap.Bool("order", group.Ordered))
			rw, err := dataChannel.Detach()
			if err != nil {
				handler.Error("failed to detach data channel", zap.String("name", group.GroupName), zap.Bool("order", group.Ordered), zap.Error(err))
				return
			}

			conn := NewWebRTCConn(peer, dataChannel, rw, &ConnInfo{
				AcceptorType:       handler.config.Type,
				AcceptorName:       handler.config.Name,
				AcceptorGroup:      group.GroupName,
				Ordered:            group.Ordered,
				BucketFillInterval: group.BucketFillInterval,
				BucketCapacity:     group.BucketCapacity,
			})
			conn.remoteAddr = ReadUserIP(r)
			conn.localAddr = handler.localAddr

			handler.connChan <- conn

			handler.Info("create datachannel success",
				zap.String("name", handler.config.Name),
				zap.String("group", group.GroupName),
				zap.String("host", r.Host),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("referer", r.Referer()),
				zap.String("uri", r.RequestURI))

		})

	}
}

func ReadUserIP(r *http.Request) *net.TCPAddr {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	portStr := strings.Split(IPAddress, ":")[1]
	port := 0
	port, _ = strconv.Atoi(portStr)
	return &net.TCPAddr{
		IP:   net.ParseIP(strings.Split(IPAddress, ":")[0]),
		Port: port,
	}
}
