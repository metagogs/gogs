package codec

import (
	"encoding/json"
)

type JSONEncode struct{}

func (c *JSONEncode) Encode(in interface{}) ([]byte, error) {
	return json.Marshal(in)
}

func (c *JSONEncode) String(in interface{}) string {
	msg, err := json.Marshal(in)
	if err != nil {
		return "proto encode error, " + err.Error()
	}

	return string(msg)
}
