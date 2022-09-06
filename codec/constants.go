package codec

import "errors"

var (
	ErrActionNotFound         = errors.New("action not found")
	ErrMethodNotFound         = errors.New("method not found")
	ErrMessageDecode          = errors.New("message decode error")
	ErrRegisterCodecType      = errors.New("register codec type error")
	ErrRegisterCodecTypeExist = errors.New("register codec type exist")
	ErrCodecType              = errors.New("codec type error")
	ErrNotValidJSONType       = errors.New("not valid json type")
	ErrActionNotExist         = errors.New("action not exist")
	ErrInvalidProtoMessage    = errors.New("invalid proto message")
)

const (
	CodecJSONDataNoHeader = uint8(0) // 标准协议头的JSON
	CodecJSONData         = uint8(1) // 标准协议头的JSON
	CodecProtoData        = uint8(2) // 使用proto解码

	CurrentMaxCodecType = uint8(2) // 当前最大编码类型
	MaxCodecType        = uint8(7) // 最大编码类型
)
