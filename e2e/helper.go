package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/metagogs/gogs/config"
	"github.com/metagogs/gogs/e2e/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func startServer(t *testing.T, config *config.Config) (context.CancelFunc, chan struct{}) { //nolint
	defer func() {
		if err := recover(); err != nil {
			t.Fatal("start error", err)
		}
	}()
	startTest := make(chan struct{})
	config.StaredCallback = func() {
		<-time.After(1 * time.Second)
		startTest <- struct{}{}
	}
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	go testdata.StartServer(ctx, config)
	return cancel, startTest
}

type TestClient struct {
	*websocket.Conn
	Datas chan []byte
}

func (c *TestClient) Start(t *testing.T) {
	t.Helper()
	for {
		_, data, err := c.ReadMessage()
		if err != nil {
			break
		}
		c.Datas <- data

	}
}
func (c *TestClient) Start2() {
	for {
		_, data, err := c.ReadMessage()
		if err != nil {
			break
		}
		c.Datas <- data

	}
}

func NewWSClinet(address string) (*TestClient, error) { //nolint
	conn, _, err := websocket.DefaultDialer.Dial(address, nil) //nolint
	if err != nil {
		return nil, err
	}

	return &TestClient{
		Conn:  conn,
		Datas: make(chan []byte),
	}, nil

}

func runningInfo(t *testing.T) *gjson.Result { //nolint
	t.Helper()
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://127.0.0.1:9999/admin", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	result := gjson.ParseBytes(bytes)
	return &result
}

func userLogin(t *testing.T, name string) string { //nolint
	t.Helper()
	data, _ := json.Marshal(map[string]string{ //nolint
		"username": name,
		"password": "e2epassword",
	})
	bf := bytes.NewReader(data)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "http://127.0.0.1:8890/user/login", bf)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	result := gjson.ParseBytes(bytes)
	return result.Get("data").Get("uid").String()
}

func encodeMessage(t *testing.T, in interface{}) []byte { //nolint
	t.Helper()
	p, err := testdata.TestApp.MessageServer.EncodeMessage(in)
	assert.Nil(t, err)
	return p.ToData().B
}
