package codec

import (
	"reflect"

	"github.com/gogf/gf/util/gconv"
	"github.com/metagogs/gogs/config"
	"github.com/metagogs/gogs/gslog"
	"github.com/metagogs/gogs/packet"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"

	"github.com/bytedance/sonic"
)

// 7E 0A 41 00 02 00 00 0F 0A 0D 31 36 36 30 33 36 30 39 31 35 35 35 36

// Dispatch 用于获取对象类型和ActionKey
type Dispatch interface {
	GetAction(in interface{}) (uint32, error) //直接通过反射获取对象的ActionKey
	GetActionByName(string) (uint32, error)   //通过名称获取对象的ActionKey
	GetObj(uint32) (interface{}, bool)
}

type PacketDecode interface {
	Decode(data []byte, in interface{}) error
	String(in interface{}) string
}

type PacketEncode interface {
	Encode(in interface{}) ([]byte, error)
	String(in interface{}) string
}

type CodecHelper struct {
	config           *config.Config
	dispatchServer   Dispatch               //用于获取对象类型和ActionKey
	packetDecodes    map[uint8]PacketDecode //解码器
	packetDecodeName map[string]uint8       //解码器名称
	packetEncode     PacketEncode           //编码器
	packetEncodeType uint8                  //编码器类型
	packetVersion    uint8                  //协议版本
	*zap.Logger
}

func NewCodecHelper(config *config.Config, dispatcher Dispatch) *CodecHelper {
	helper := &CodecHelper{
		config:           config,
		dispatchServer:   dispatcher,
		packetDecodes:    make(map[uint8]PacketDecode),
		packetDecodeName: make(map[string]uint8),
		packetVersion:    1,
		Logger:           gslog.NewLog("codec"),
	}

	//注册解码器
	protoDecode := &ProtoDecode{}
	_ = helper.RegisterDecode(CodecProtoData, protoDecode)

	jsonDecode := &JSONDecode{}
	_ = helper.RegisterDecode(CodecJSONData, jsonDecode)

	//注册编码器
	protoEncode := &ProtoEncode{}
	helper.RegisterEncode(CodecProtoData, protoEncode)

	return helper
}

func (d *CodecHelper) RegisterDecode(decodeKey uint8, c PacketDecode) error {
	if decodeKey > MaxCodecType {
		d.Warn("register decode type error max type", zap.String("packet_decode", reflect.TypeOf(c).Elem().Name()))
		return ErrRegisterCodecType
	}
	if _, ok := d.packetDecodes[decodeKey]; ok {
		d.Warn("register decode type error exist", zap.String("packet_decode", reflect.TypeOf(c).Elem().Name()))
		return ErrRegisterCodecTypeExist
	}

	d.packetDecodes[decodeKey] = c
	d.packetDecodeName[reflect.TypeOf(c).Elem().Name()] = decodeKey
	d.Info("register decode type success", zap.String("packet_decode", reflect.TypeOf(c).Elem().Name()))

	return nil
}

func (d *CodecHelper) RegisterEncode(v uint8, e PacketEncode) {
	d.packetEncode = e
	d.packetEncodeType = v
	d.Info("register encode type success", zap.String("packet_encode", reflect.TypeOf(e).Elem().Name()))
}

func (d *CodecHelper) GetTypes() (map[string]uint8, string) {
	return d.packetDecodeName, reflect.TypeOf(d.packetEncode).Elem().Name()
}

func (d *CodecHelper) Decode(data []byte) (*packet.Packet, error) {
	parsePacket, err := packet.ParsePacket(data)
	if err != nil {
		//尝试判断是不是文本类型或者二进制的json类型
		p, e := d.decodeJSONWithoutHeader(data)
		if e != nil {
			//如果不是特殊类型，则按照协议错误处理
			d.Warn("parse packet error", zap.Error(err), zap.String("json_error", e.Error()))
			return nil, err
		}

		if d.config.Debug || d.config.ReceiveMessageLog {
			//调试模式，打印所有信息
			d.Info(string(p.Data),
				zap.String("type", "json_without_header"),
				zap.String("codec", "decode"))
		}

		return p, nil
	}

	if d.config.Debug || d.config.ReceiveMessageLog {
		d.Info("codec get new packet",
			zap.Any("packet_type", parsePacket.GetPacketType()),
			zap.Uint8("encode_type", parsePacket.GetEncodeType()),
			zap.Uint8("version", parsePacket.GetVersion()),
			zap.Uint8("module", parsePacket.GetModule()),
			zap.Uint32("action", parsePacket.GetActionKey()),
			zap.Uint32("length", parsePacket.GetLength()))
	}

	if helper, ok := d.packetDecodes[parsePacket.GetEncodeType()]; ok {

		actionKey := parsePacket.GetActionKey()
		in, ok := d.dispatchServer.GetObj(actionKey)
		if !ok {
			return nil, ErrActionNotFound
		}

		if err := helper.Decode(parsePacket.Data, in); err != nil {
			d.Info("decode error", zap.Error(err))
			return nil, err
		}
		parsePacket.Bean = in

		if d.config.Debug || d.config.ReceiveMessageLog {
			//调试模式，打印所有信息
			d.Info(helper.String(in),
				zap.String("codec", "decode"),
				zap.String("type", reflect.TypeOf(helper).Elem().Name()),
				zap.Uint32("action_key", actionKey))
		}

		return parsePacket, nil
	}

	d.Info("decode type not found", zap.Any("packet_type", parsePacket.GetPacketType()))
	return nil, ErrCodecType
}

func (d *CodecHelper) decodeJSONWithoutHeader(data []byte) (*packet.Packet, error) {
	if len(data) < 2 {
		//一定不是个JSON，没必要了
		return nil, ErrNotValidJSONType
	}
	if gjson.ValidBytes(data) {
		actionName := gjson.ParseBytes(data).Get("action").String()
		if len(actionName) == 0 {
			return nil, ErrActionNotExist
		}

		actionKey, err := d.dispatchServer.GetActionByName(actionName)
		if err != nil {
			return nil, err
		}

		in, ok := d.dispatchServer.GetObj(actionKey)
		if !ok {
			return nil, ErrActionNotExist
		}

		if err := sonic.Unmarshal(data, in); err != nil {
			return nil, err
		}

		//重新构建一个header和packet
		jsonPacket := packet.NewPacketWithHeader(data, d.packetVersion, CodecJSONDataNoHeader, actionKey)
		jsonPacket.Bean = in

		return jsonPacket, nil
	}

	return nil, ErrNotValidJSONType
}

func (d *CodecHelper) Encode(in interface{}, name ...string) (*packet.Packet, error) {
	if d.packetEncodeType == CodecJSONDataNoHeader {
		return d.encodeJSONWithoutHeaderEncode(in, name...)
	}

	var action uint32
	var err error
	if len(name) > 0 && len(name[0]) > 0 {
		action, _ = d.dispatchServer.GetActionByName(name[0])
	}
	if action == 0 {
		action, err = d.dispatchServer.GetAction(in)
		if err != nil {
			return nil, err
		}
	}

	data, err := d.packetEncode.Encode(in)
	if err != nil {
		return nil, err
	}

	packetData := packet.NewPacketWithHeader(data, d.packetVersion, d.packetEncodeType, action)

	if d.config.Debug || d.config.SendMessageLog {
		d.Info(d.packetEncode.String(in),
			zap.String("codec", "decode"),
			zap.String("type", reflect.TypeOf(d.packetEncode).Elem().Name()),
			zap.Uint32("action", action))
	}

	return packetData, nil
}

func (d *CodecHelper) encodeJSONWithoutHeaderEncode(in interface{}, name ...string) (*packet.Packet, error) {
	inMap := gconv.Map(in)
	if inMap == nil {
		inMap = make(map[string]interface{})
	}

	if len(name) > 0 && len(name[0]) > 0 {
		inMap["action"] = name[0]
	} else {
		inType := reflect.TypeOf(in)
		if inType.Kind() == reflect.Ptr {
			inType = inType.Elem()
		}
		inMap["action"] = inType.Name()
	}

	data, err := d.packetEncode.Encode(inMap)
	if err != nil {
		return nil, err
	}

	if d.config.Debug {
		d.Info(string(data), zap.String("codec", "encode"), zap.String("type", "json_without_header"))
	}

	return packet.NewPacket(data), nil
}
