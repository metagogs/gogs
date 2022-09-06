package codec

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type ProtoEncode struct{}

func (c *ProtoEncode) Encode(in interface{}) ([]byte, error) {
	pb, ok := in.(proto.Message)
	if !ok {
		return nil, ErrInvalidProtoMessage
	}
	return proto.Marshal(pb)
}

func (c *ProtoEncode) String(in interface{}) string {
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
