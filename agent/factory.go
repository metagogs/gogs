package agent

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/metagogs/gogs/acceptor"
	"github.com/metagogs/gogs/config"
	"github.com/metagogs/gogs/gslog"
	"github.com/metagogs/gogs/message"
	"github.com/metagogs/gogs/proto"
	"github.com/metagogs/gogs/session"
	"github.com/metagogs/gogs/utils/bytebuffer"
	"github.com/metagogs/gogs/utils/snow"
	"go.uber.org/zap"
)

var (
	pingPool = &sync.Pool{
		New: func() interface{} {
			return &proto.Ping{}
		},
	}
)

// AgentFactory is the factory to create the agent
type AgentFactory struct {
	config             *config.Config
	sf                 *snowflake.Node        // snowflake node
	sessionPool        session.SessionPool    // session pool
	messageServer      *message.MessageServer // message server
	messagesBufferSize int                    // message buffer size
}

func NewAgentFactory(config *config.Config,
	pool session.SessionPool,
	messageServer *message.MessageServer) *AgentFactory {

	sf, err := snow.NewSnowNode()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return &AgentFactory{
		config:             config,
		sf:                 sf,
		sessionPool:        pool,
		messageServer:      messageServer,
		messagesBufferSize: config.AgentMessageBufferSize,
	}
}

func (af *AgentFactory) NewAgent(conn acceptor.AcceptorConn) *Agent {
	agentId := af.sf.Generate().Int64()
	agent := &Agent{
		AgentID:          agentId,
		conn:             conn,
		agentLog:         gslog.NewLog("agent").With(zap.Int64("agent_id", agentId)),
		chSend:           make(chan *bytebuffer.ByteBuffer, af.messagesBufferSize),
		chSendByte:       make(chan []byte, af.messagesBufferSize),
		chStopWrite:      make(chan struct{}),
		chStopHeartbeat:  make(chan struct{}),
		chDie:            make(chan struct{}),
		messageServer:    af.messageServer,
		heartbeatTimeout: time.Duration(af.config.AgentHeartBeatTimeout) * time.Second,
		heartbeatLog:     af.config.AgentHeartBeatLog,
	}
	// create session
	s := af.sessionPool.CreateSession(agent)
	agent.sess = s

	return agent
}
