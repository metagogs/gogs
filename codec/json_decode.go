package codec

import (
	"github.com/bytedance/sonic"
)

type JSONDecode struct{}

func (c *JSONDecode) Decode(data []byte, in interface{}) error {
	if sonic.ConfigDefault.Valid(data) {
		return sonic.Unmarshal(data, in)
	}
	return ErrNotValidJSONType
}

func (c *JSONDecode) String(in interface{}) string {
	msg, err := sonic.Marshal(in)
	if err != nil {
		return "proto encode error, " + err.Error()
	}

	return string(msg)
}
