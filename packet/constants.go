package packet

import "errors"

const (
	// 协议头长度
	HeaderLength = 8
	// 协议头标识
	HeaderFlag = 0x7E
	// 最大内容程度，16MBB
	MaxPacketSize = 1 << 24
)

type PacketType byte

const (
	// 系统内置包，如Ping Pong
	SystemPacket PacketType = 0x01
	// 服务类型包，如Login、Register
	ServicePacket PacketType = 0x02
)

var ErrHeaderFlag = errors.New("wrong header flag")
var ErrHeaderLength = errors.New("wrong header length")
var ErrMaxPacketSize = errors.New("max packet size")

var ErrEncodeType = errors.New("wrong encode type")
