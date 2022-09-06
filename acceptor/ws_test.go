package acceptor

import (
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/metagogs/gogs/global"
	"github.com/stretchr/testify/assert"
)

func TestWebSockets(t *testing.T) {
	global.GOGS_DISABLE_LOG = true
	acceptor := NewWSAcceptror(&AcceptroConfig{
		Name:     "websockets",
		HttpPort: 11003,
		Groups: []*AcceptorGroupConfig{
			{
				GroupName: "testbase",
			},
		},
	})
	go func() {
		acceptor.ListenAndServe()
	}()
	<-time.After(2 * time.Second)

	client, err := newWSClinet("ws://127.0.0.1:11003/testbase")
	assert.Nil(t, err)
	go client.Start(t)

	var serverConn AcceptorConn
	select {
	case <-time.After(1 * time.Second):
		t.Fatal("get acceptor conn timeout")
	case serverConn = <-acceptor.GetConnChan():
		t.Log("get acceptor conn success")
	}

	assert.Equal(t, "testbase", serverConn.GetInfo().AcceptorGroup)
	assert.Equal(t, "websockets", serverConn.GetInfo().AcceptorName)
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
	case msg := <-client.datas:
		t.Log("client receive message", string(msg))
		assert.EqualValues(t, []byte("hello"), msg)
	}

	t.Log("client send message")
	err = client.WriteMessage(websocket.BinaryMessage, []byte("world"))
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

func newWSClinet(address string) (*client, error) {
	conn, _, err := websocket.DefaultDialer.Dial(address, nil)
	if err != nil {
		return nil, err
	}

	return &client{
		Conn:  conn,
		datas: make(chan []byte),
	}, nil

}

type client struct {
	*websocket.Conn
	datas chan []byte
}

func (c *client) Start(t *testing.T) {
	t.Helper()
	for {
		_, data, err := c.ReadMessage()
		if err != nil {
			break
		}
		c.datas <- data

	}
}
