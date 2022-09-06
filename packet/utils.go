package packet

func CreateHeader(data []byte, version, encodeType uint8, action uint32) []byte {
	headerByte := make([]byte, 8)
	headerByte[0] = HeaderFlag
	headerByte[1] = version<<3 | encodeType
	headerByte[2] = byte(action >> 16)
	headerByte[3] = byte(action >> 8)
	headerByte[4] = byte(action)
	length := len(data)
	headerByte[5] = byte(length >> 16)
	headerByte[6] = byte(length >> 8)
	headerByte[7] = byte(length)

	return headerByte
}

func CreateAction(packetType PacketType, module uint8, action uint16) uint32 {
	return uint32(packetType)<<22 | uint32(module)<<16 | uint32(action)
}

func ActionToBytes(action uint32) []byte {
	return []byte{
		byte(action >> 16),
		byte(action >> 8),
		byte(action),
	}
}
