package message

import (
	"context"

	"github.com/metagogs/gogs/codec"
	"github.com/metagogs/gogs/component"
	"github.com/metagogs/gogs/dispatch"
	"github.com/metagogs/gogs/packet"
	"github.com/metagogs/gogs/session"
)

type MessageServer struct {
	codecHelper    *codec.CodecHelper
	dispatchServer *dispatch.DispatchServer
}

func NewMessageServer(codecHelper *codec.CodecHelper, dispatchServer *dispatch.DispatchServer) *MessageServer {
	return &MessageServer{
		codecHelper:    codecHelper,
		dispatchServer: dispatchServer,
	}
}

// DecodeMessage 解码消息
func (m *MessageServer) DecodeMessage(data []byte) (*packet.Packet, error) {
	return m.codecHelper.Decode(data)
}

// CallMessageHandler 调用消息处理器
func (m *MessageServer) CallMessageHandler(ctx context.Context, sess *session.Session, in *packet.Packet) error {
	return m.dispatchServer.Call(ctx, sess, in)
}

// EncodeMessage 编码消息
func (m *MessageServer) EncodeMessage(in interface{}) (*packet.Packet, error) {
	return m.codecHelper.Encode(in)
}

// RegisterComponent 注册组件
func (m *MessageServer) RegisterComponent(sd component.ComponentDesc, ss interface{}) {
	m.dispatchServer.RegisterComponent(sd, ss)
}

// UseDefaultEncodeJSON 设置默认编码器
func (m *MessageServer) UseDefaultEncodeJSON() {
	m.codecHelper.RegisterEncode(codec.CodecJSONData, &codec.JSONEncode{})
}

// UseDefaultEncodeJSON 设置默认编码器
func (m *MessageServer) UseDefaultEncodeProto() {
	m.codecHelper.RegisterEncode(codec.CodecProtoData, &codec.ProtoEncode{})
}

// UseDefaultEncodeJSON 设置默认编码器
func (m *MessageServer) UseDefaultEncodeJSONWithHeader() {
	// 发送消息时不包含头部,使用JSON,方便测试使用
	m.codecHelper.RegisterEncode(codec.CodecJSONDataNoHeader, &codec.JSONEncode{})
}
