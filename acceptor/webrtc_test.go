package acceptor

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/metagogs/gogs/global"
	"github.com/pion/webrtc/v3"
	"github.com/stretchr/testify/assert"
)

func TestWebRTC(t *testing.T) {
	global.GOGS_DISABLE_LOG = true
	acceptor := NewWebRTCAcceptor(&AcceptroConfig{
		Name:     "wertc",
		HttpPort: 11001,
		UdpPort:  11002,
		Groups: []*AcceptorGroupConfig{
			{
				GroupName: "testchannel",
			},
		},
	})
	go func() {
		acceptor.ListenAndServe()
	}()
	<-time.After(2 * time.Second)

	connectSuccess := make(chan bool)
	clientMessage := make(chan []byte)
	var clientChannel *webrtc.DataChannel

	// create webrtc connection by client
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	_, err = peerConnection.CreateDataChannel("testchannel", nil)
	assert.Nil(t, err)
	peerConnection.OnDataChannel(func(dc *webrtc.DataChannel) {
		t.Log("data channel open")
		connectSuccess <- true
		clientChannel = dc
		clientChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
			t.Log("get message from", dc.Label())
			clientMessage <- msg.Data
		})
	})

	assert.Nil(t, err)
	offer, err := peerConnection.CreateOffer(nil)
	assert.Nil(t, err)
	err = peerConnection.SetLocalDescription(offer)
	assert.Nil(t, err)

	payload, err := json.Marshal(offer)
	assert.Nil(t, err)
	resp, err := http.Post("http://127.0.0.1:11001", "application/json", bytes.NewReader(payload))
	defer resp.Body.Close()
	assert.Nil(t, err)

	sdp := webrtc.SessionDescription{}
	err = json.NewDecoder(resp.Body).Decode(&sdp)
	assert.Nil(t, err)
	err = peerConnection.SetRemoteDescription(sdp)
	assert.Nil(t, err)

	select {
	case <-time.After(1 * time.Second):
		t.Fatal("client connect server timeout")
	case <-connectSuccess:
		t.Log("connect server success")
	}

	var serverConn AcceptorConn
	select {
	case <-time.After(1 * time.Second):
		t.Fatal("get acceptor conn timeout")
	case serverConn = <-acceptor.GetConnChan():
		t.Log("get acceptor conn success")
	}

	assert.Equal(t, "testchannel", serverConn.GetInfo().AcceptorGroup)
	assert.Equal(t, "wertc", serverConn.GetInfo().AcceptorName)
	serverMsg := make(chan []byte)
	go func() {
		for {
			msg, err := serverConn.GetNextMessage()
			if err != nil {
				return
			}
			serverMsg <- msg
		}
	}()

	t.Log("server send message")
	n, err := serverConn.Write([]byte("hello"))
	assert.Equal(t, len([]byte("hello")), n)
	assert.Nil(t, err)

	select {
	case <-time.After(1 * time.Second):
		t.Fatal("client receive message timeout")
	case msg := <-clientMessage:
		t.Log("client receive message", string(msg))
		assert.EqualValues(t, []byte("hello"), msg)
	}

	t.Log("client send message")
	err = clientChannel.Send([]byte("world"))
	assert.Nil(t, err)

	select {
	case <-time.After(1 * time.Second):
		t.Fatal("server receive message timeout")
	case msg := <-serverMsg:
		t.Log("server receive message", string(msg))
		assert.Nil(t, err)
		assert.EqualValues(t, []byte("world"), msg)
	}

	err = serverConn.Close()
	assert.Nil(t, err)

	closed := serverConn.IsClosed()
	assert.True(t, closed)

	acceptor.Stop()

	<-time.After(2 * time.Second)

}
