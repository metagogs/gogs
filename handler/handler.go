package handler

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"math"
	"runtime/debug"

	"github.com/metagogs/gogs/acceptor"
	"github.com/metagogs/gogs/agent"
	"github.com/metagogs/gogs/config"
	"github.com/metagogs/gogs/gslog"
	"github.com/metagogs/gogs/message"
	"github.com/metagogs/gogs/packet"
	"github.com/metagogs/gogs/session"
	"go.uber.org/zap"
)

type unHandlerMessage struct {
	ctx  context.Context
	sess *session.Session
	in   *packet.Packet
}

type HandlerService struct {
	config                *config.Config
	agentFactory          *agent.AgentFactory
	messageServer         *message.MessageServer
	chLocalMessageProcess chan unHandlerMessage
	*zap.Logger
}

// NewHandlerService create a new handler service, can get all the client's message
func NewHandlerService(config *config.Config, factory *agent.AgentFactory, messageServer *message.MessageServer) *HandlerService {
	return &HandlerService{
		config:                config,
		agentFactory:          factory,
		Logger:                gslog.NewLog("handler"),
		messageServer:         messageServer,
		chLocalMessageProcess: make(chan unHandlerMessage, 100),
	}
}

func (h *HandlerService) SetAgentFactory(factory *agent.AgentFactory) {
	h.agentFactory = factory
}

func (h *HandlerService) dispatchMessage(m unHandlerMessage) {
	defer func() {
		if r := recover(); r != nil {
			h.toPanicError(r)
			h.Error("dispatch message error", zap.Any("recover", r))
		}
	}()
	if err := h.messageServer.CallMessageHandler(m.ctx, m.sess, m.in); err != nil {
		h.Error("dispatch message error",
			zap.Error(err),
			zap.Int64("agent_id", m.sess.ID()))
	}
}

func (h *HandlerService) toPanicError(r interface{}) {
	buf := bytes.Buffer{}
	stackScanner := bufio.NewScanner(bytes.NewReader(debug.Stack()))
	for i := 0; i < math.MaxInt32; i++ {
		if !stackScanner.Scan() {
			break
		}

		text := stackScanner.Text()
		buf.WriteString(text + "\n")
	}
	h.Error(buf.String())
}

func (h *HandlerService) Handle(conn acceptor.AcceptorConn) {
	a := h.agentFactory.NewAgent(conn)
	h.Info("get new agent connection in handle", zap.Int64("agent_id", a.AgentID))

	go a.Start() // start agent

	defer func() {
		a.GetSession().Close()
	}()

	for {
		if a.IsClosed() {
			return
		}
		msg, err := conn.GetNextMessage()
		if errors.Is(err, acceptor.ErrMessageRateLimit) {
			continue
		}
		if err != nil {
			h.Warn("get next message error", zap.Error(err), zap.Int64("agent_id", a.AgentID))
			return
		}
		if msg == nil {
			continue
		}

		h.processPacket(a, msg)
	}
}

func (h *HandlerService) processPacket(a *agent.Agent, data []byte) {
	if h.config.ReceiveMessageLog {
		h.Info("agent get new packet",
			zap.Int64("agent_id", a.AgentID),
			zap.String("remote_addr", a.RemoteAddr().String()))
	}

	in, err := h.messageServer.DecodeMessage(data)
	if err != nil {
		h.Error("decode packet error", zap.Error(err), zap.Int64("agent_id", a.AgentID))
		return
	}

	h.dispatchMessage(unHandlerMessage{
		ctx:  context.Background(),
		sess: a.GetSession(),
		in:   in,
	})

	a.SetLastAt()
}
