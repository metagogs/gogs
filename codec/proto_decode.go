package codec

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type ProtoDecode struct{}

func (c *ProtoDecode) Decode(data []byte, in interface{}) error {
	pb, ok := in.(proto.Message)
	if !ok {
		return ErrInvalidProtoMessage
	}
	return proto.Unmarshal(data, pb)
}

func (c *ProtoDecode) String(in interface{}) string {
	pb, ok := in.(proto.Message)
	if !ok {
		return "proto encode error, is not valid proto message"
	}

	msg, err := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(pb)
	if err != nil {
		return "proto encode error, " + err.Error()
	}

	return string(msg)
}
