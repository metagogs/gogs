package codec

import (
	"encoding/json"
)

type JSONDecode struct{}

func (c *JSONDecode) Decode(data []byte, in interface{}) error {
	if json.Valid(data) {
		return json.Unmarshal(data, in)
	}
	return ErrNotValidJSONType
}

func (c *JSONDecode) String(in interface{}) string {
	msg, err := json.Marshal(in)
	if err != nil {
		return "proto encode error, " + err.Error()
	}

	return string(msg)
}
