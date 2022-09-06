package packet

import (
	"github.com/metagogs/gogs/utils/bytebuffer"
)

// package header 64bit

// 0byte      flag: 0x7E  8bit  协议头固定

// 1byte      version:    5bit  协议版本
//            encodeType: 3bit  协议编码类型

// 2byte      packetType: 2bit  消息类型 0x01: system, 0x02: service
//            module:     6bit  消息模块

// 3.4byte    action:     16bit 消息ID

// 5.6.7byte  length:     24bit 消息长度
type Packet struct {
	Bean       interface{}
	HeaderByte []byte
	Data       []byte
}

// ParsePacket 从字节数组中解析出Packet
func ParsePacket(data []byte) (*Packet, error) {
	dataLength := len(data)
	if dataLength < HeaderLength {
		return nil, ErrHeaderLength
	}
	if dataLength > MaxPacketSize {
		return nil, ErrMaxPacketSize
	}

	header := data[:HeaderLength]
	//判断标识符
	if header[0] != HeaderFlag {
		return nil, ErrHeaderFlag
	}

	packetType := PacketType(header[2] >> 6)
	if packetType != SystemPacket && packetType != ServicePacket {
		return nil, ErrEncodeType
	}

	packet := GetPacket()
	packet.HeaderByte = header
	packet.Data = data[HeaderLength:]

	return packet, nil
}

func NewPacket(data []byte) *Packet {
	packet := GetPacket()
	packet.Data = data
	packet.HeaderByte[0] = 0x20

	return packet
}

func NewPacketWithHeader(data []byte, version, encodeType uint8, action uint32) *Packet {
	packet := GetPacket()
	packet.HeaderByte[0] = HeaderFlag
	packet.HeaderByte[1] = version<<3 | encodeType
	packet.HeaderByte[2] = byte(action >> 16)
	packet.HeaderByte[3] = byte(action >> 8)
	packet.HeaderByte[4] = byte(action)
	length := len(data)
	packet.HeaderByte[5] = byte(length >> 16)
	packet.HeaderByte[6] = byte(length >> 8)
	packet.HeaderByte[7] = byte(length)
	packet.Data = data

	return packet
}

func (p *Packet) ToData() *bytebuffer.ByteBuffer {
	defer func() {
		PutPacket(p)
	}()

	b := bytebuffer.Get()
	if p.HeaderByte[0] == HeaderFlag {
		_, _ = b.Write(p.HeaderByte)
	}
	_, _ = b.Write(p.Data)

	return b
}

// packetType: 2bit  // 0x01: system, 0x02: service
// module:     6bit
// action:     16bit
// 合一起的值 协议头/2/3/4字节
func (p *Packet) GetActionKey() uint32 {
	return uint32(p.HeaderByte[2])<<16 | uint32(p.HeaderByte[3])<<8 | uint32(p.HeaderByte[4])
}

func (p *Packet) GetPacketType() PacketType {
	return PacketType(p.HeaderByte[2] >> 6)
}

func (p *Packet) GetVersion() uint8 {
	return p.HeaderByte[1] >> 3
}

func (p *Packet) GetEncodeType() uint8 {
	return p.HeaderByte[1] & (0xff >> 5)
}

func (p *Packet) GetModule() uint8 {
	return p.HeaderByte[2] & (0xff >> 6)
}

func (p *Packet) GetLength() uint32 {
	return uint32(p.HeaderByte[5])<<16 | uint32(p.HeaderByte[6])<<8 | uint32(p.HeaderByte[7])
}
