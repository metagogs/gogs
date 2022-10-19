package message

var DefaultMessageServer *MessageServer

func EncodeMessage(in interface{}, name ...string) ([]byte, error) {
	packet, err := DefaultMessageServer.EncodeMessage(in, name...)
	if err != nil {
		return nil, err
	}
	return packet.ToByte(), nil
}
