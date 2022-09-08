package e2e

import (
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/metagogs/gogs/config"
	"github.com/metagogs/gogs/e2e/testdata"
	"github.com/metagogs/gogs/e2e/testdata/fakeinternal/logic/baseworld"
	"github.com/metagogs/gogs/e2e/testdata/game"
	"github.com/metagogs/gogs/global"
	"github.com/metagogs/gogs/session"
	"github.com/stretchr/testify/assert"
)

var (
	defaultConfig = config.NewConfig("default.yaml")
	testClient    *client
	testClient2   *client
	uid           string
	uid2          string
)

// TestSendConnectWS test the base function
// client connect to the server
// client send the message to the server
// server sned the message to the client
// we should check the message can be received and the message is equal to the message we send
// todo add the webrtc datachanel test in the future
func TestSendConnectWS(t *testing.T) {
	global.GOGS_DISABLE_LOG = true
	global.GoGSDebug = true
	// start the gogs server
	defaultConfig.AgentHeartBeatTimeout = 0
	cancel, started := startServer(t, defaultConfig)
	defer cancel()
	// wait for started
	<-started
	t.Log("websocket testing")
	clients := []*client{}
	// try to make 10 websocket client
	for i := 0; i < 10; i++ {
		client, err := newWSClinet("ws://127.0.0.1:8888/base")
		assert.Nil(t, err)
		clients = append(clients, client)
		defer client.Close()
	}

	// check the running info, if the connection count is right to 10
	info := runningInfo(t)
	currentConnections := info.Get("online_session_count").String()
	assert.Equal(t, "10", currentConnections)

	// close some client
	for i := 0; i < 3; i++ {
		clients[i].Close()
	}
	<-time.After(1 * time.Second)
	// check the running info again
	info = runningInfo(t)
	currentConnections = info.Get("online_session_count").String()
	assert.Equal(t, "7", currentConnections)

	totalConnections := info.Get("session_count").String()
	assert.Equal(t, "10", totalConnections)

	// close all the client
	for i := 3; i < 10; i++ {
		clients[i].Close()
	}

	var err error
	t.Log("websocket message testing")
	t.Log("user1 login")
	testClient, err = newWSClinet("ws://127.0.0.1:8888/base")
	assert.Nil(t, err)
	go testClient.Start(t)

	<-time.After(1 * time.Second)

	// the session pool should only have one session
	sessions := session.ListSessions()
	assert.Equal(t, 1, len(sessions))

	//test send data
	sessions[0].SendData([]byte("hello world"))
	select {
	case <-time.After(1 * time.Second):
		t.Fatal("server get the data timeout")
	case msg := <-testClient.datas:
		assert.Equal(t, "hello world", string(msg))
	}

	uid = userLogin(t, "e2e")
	assert.NotEmpty(t, uid)

	bindUser(t, testClient, uid)
	joinWorld(t, testClient)

	t.Log("user2 login")
	testClient2, err = newWSClinet("ws://127.0.0.1:8888/base")
	assert.Nil(t, err)
	go testClient2.Start(t)

	<-time.After(1 * time.Second)

	// with the two users logined, we should have to sessions
	sessions = session.ListSessions()
	assert.Equal(t, 2, len(sessions))

	uid2 = userLogin(t, "e2e2")
	assert.NotEmpty(t, uid2)

	t.Log("use the json encode as the encode type")
	testdata.TestApp.UseDefaultEncodeJSON()
	bindUser(t, testClient2, uid2)
	joinWorld(t, testClient2)

	//client1 get message
	select {
	case <-time.After(1 * time.Second):
		t.Fatal("client get the JoinWorldNotify message timeout")
	case msg := <-testClient.datas:
		t.Log("client get the JoinWorldNotify message")
		p, err := testdata.TestApp.MessageServer.DecodeMessage(msg)
		assert.Nil(t, err)
		result, ok := p.Bean.(*game.JoinWorldNotify)
		assert.True(t, ok)
		assert.Equal(t, uid2, result.Uid)
		assert.Equal(t, "e2e2", result.Name)
	}

	t.Log("change the encode type with json encode without header")
	testdata.TestApp.UseDefaultEncodeJSONWithHeader()
	t.Log("user1 send the UpdateUserInWorld message")
	err = testClient.WriteMessage(websocket.BinaryMessage, encodeMessage(t, &game.UpdateUserInWorld{
		Uid: uid,
		Position: &game.Vecotr3{
			X: 11,
			Y: 22,
			Z: 33,
		},
	}))
	assert.Nil(t, err)

	select {
	case <-time.After(1 * time.Second):
		t.Fatal("client get t he UpdateUserInWorld message timeout")
	case msg := <-testClient2.datas:
		t.Log("client get the UpdateUserInWorld message")
		p, err := testdata.TestApp.MessageServer.DecodeMessage(msg)
		assert.Nil(t, err)
		result, ok := p.Bean.(*game.UpdateUserInWorld)
		assert.True(t, ok)
		assert.Equal(t, uid, result.Uid)
		assert.Equal(t, float32(11), result.Position.X)
		assert.Equal(t, float32(22), result.Position.Y)
		assert.Equal(t, float32(33), result.Position.Z)
	}

	endTime := time.After(2 * time.Second)
	<-endTime
}

// bindUser when we create a client, we should send the BindUser message to th server
// to bind the user to the connection. The servet should send the BindUserSuccess message to the client
func bindUser(t *testing.T, sendClient *client, id string) {
	t.Helper()
	t.Log("client send BindUser message")
	err := sendClient.WriteMessage(websocket.BinaryMessage, encodeMessage(t, &game.BindUser{
		Uid: id,
	}))
	assert.Nil(t, err)
	select {
	case <-time.After(1 * time.Second):
		t.Fatal("server get the BindUserResult message timeout")
	case msg := <-baseworld.BindUserHandler:
		t.Log("server get the BindUserResult message success")
		assert.Equal(t, id, msg.Uid)
	}

	select {
	case <-time.After(1 * time.Second):
		t.Fatal("client get the BindSuccess message timeout")
	case msg := <-sendClient.datas:
		t.Log("client get the BindSuccess message success")
		p, err := testdata.TestApp.MessageServer.DecodeMessage(msg)
		assert.Nil(t, err)
		_, ok := p.Bean.(*game.BindSuccess)
		assert.True(t, ok)
	}
}

// joinWorld user can send the JoinWorld message to the server to join the world
// the server should send the JoinWorldSuccess message to the client
func joinWorld(t *testing.T, sendClient *client) {
	t.Helper()
	t.Log("client send JoinWorld message")
	err := sendClient.WriteMessage(websocket.BinaryMessage, encodeMessage(t, &game.JoinWorld{
		Uid: uid,
	}))
	assert.Nil(t, err)
	select {
	case <-time.After(1 * time.Second):
		t.Fatal("server get the JoinWorld message timeout")
	case <-baseworld.JoinWorldHandler:
		t.Log("server get the JoinWorld message success")
	}

	select {
	case <-time.After(1 * time.Second):
		t.Fatal("client get the JoinWorldSuccess message timeout")
	case msg := <-sendClient.datas:
		t.Log("client get the JoinWorldSuccess message success")
		p, err := testdata.TestApp.MessageServer.DecodeMessage(msg)
		assert.Nil(t, err)
		result, ok := p.Bean.(*game.JoinWorldSuccess)
		assert.True(t, ok)
		assert.NotZero(t, len(result.Uids))
	}
}
