package codec

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/gogf/gf/util/gconv"
	"github.com/metagogs/gogs/config"
	"github.com/metagogs/gogs/packet"
	"google.golang.org/protobuf/proto"
)

// go test -v -cover -coverprofile=coverage.data .
// go tool cover -html=coverage.data -o coverage.html

type TestDispatch struct{}

func (t *TestDispatch) GetAction(in interface{}) (uint32, error) {
	if reflect.TypeOf(in).Elem().Name() == "Pong" {
		return 0, errors.New("not found")
	}
	//ping
	return 0x410001, nil
}

func (t *TestDispatch) GetActionByName(name string) (uint32, error) {
	if name == "Ping2" {
		return 0, errors.New("not found")
	}
	if name == "Pong" {
		return 0, errors.New("not found")
	}
	return 0x410001, nil
}

func (t *TestDispatch) GetObj(action uint32) (interface{}, bool) {
	if action == 0x430001 {
		return nil, false
	}
	return &Ping{}, true
}

func jsonData(in interface{}) []byte {
	inType := reflect.TypeOf(in)
	actionName := ""
	if inType.Kind() == reflect.Ptr {
		actionName = inType.Elem().Name()
	} else {
		actionName = inType.Name()
	}

	inMap := gconv.Map(in)
	if inMap == nil {
		inMap = make(map[string]interface{})
	}
	inMap["action"] = actionName
	b, _ := json.Marshal(inMap)

	return b
}

func jsonHeaderData(in interface{}) []byte {
	b, _ := sonic.Marshal(in)
	return b
}

func protoData(in interface{}) []byte {
	b, _ := proto.Marshal(in.(proto.Message))
	return b
}

func TestCodecHelper_Encode(t *testing.T) {
	dispatch := TestDispatch{}

	pingTest1 := &Ping{
		Time: "123",
	}
	pongTest := &Pong{}

	type fields struct {
		useJSONWithoutHeader bool
		encodeType           uint8
		encode               PacketEncode
	}
	type args struct {
		in interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *packet.Packet
		wantErr bool
	}{
		{
			name: "encode_json_without_header",
			fields: fields{
				useJSONWithoutHeader: true,
				encodeType:           CodecProtoData,
				encode:               &ProtoEncode{},
			},
			args: args{
				in: pingTest1,
			},
			want: &packet.Packet{
				HeaderByte: []byte{0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
				Data:       jsonData(pingTest1),
			},
		},
		{
			name: "encode_proto",
			fields: fields{
				useJSONWithoutHeader: false,
				encodeType:           CodecProtoData,
				encode:               &ProtoEncode{},
			},
			args: args{
				in: pingTest1,
			},
			want: &packet.Packet{
				HeaderByte: packet.CreateHeader(protoData(pingTest1), 1, CodecProtoData, 0x410001),
				Data:       protoData(pingTest1),
			},
		},
		{
			name: "encode_proto_err_type",
			fields: fields{
				useJSONWithoutHeader: false,
				encodeType:           CodecProtoData,
				encode:               &ProtoEncode{},
			},
			args: args{
				in: pongTest,
			},
			wantErr: true,
		},
		{
			name: "encode_json",
			fields: fields{
				useJSONWithoutHeader: false,
				encodeType:           CodecJSONData,
				encode:               &JSONEncode{},
			},
			args: args{
				in: pingTest1,
			},
			want: &packet.Packet{
				HeaderByte: packet.CreateHeader(jsonHeaderData(pingTest1), 1, CodecJSONData, 0x410001),
				Data:       jsonHeaderData(pingTest1),
			},
		},
		{
			name: "encode_json_err_type",
			fields: fields{
				useJSONWithoutHeader: false,
				encodeType:           CodecJSONData,
				encode:               &JSONEncode{},
			},
			args: args{
				in: pongTest,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := config.NewDefaultConfig()
			config.Debug = true
			d := NewCodecHelper(config, &dispatch)
			d.RegisterEncode(tt.fields.encodeType, tt.fields.encode)
			if tt.fields.useJSONWithoutHeader {
				d.RegisterEncode(CodecJSONDataNoHeader, &JSONEncode{})
			}

			got, err := d.Encode(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("CodecHelper.Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CodecHelper.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCodecHelper_Decode(t *testing.T) {
	dispatch := TestDispatch{}
	config := config.NewDefaultConfig()
	config.Debug = true
	d := NewCodecHelper(config, &dispatch)

	ping := &Ping{Time: "1660360915556"}

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *packet.Packet
		wantErr bool
	}{
		{
			name: "proto",
			args: args{
				data: append(packet.CreateHeader(protoData(ping), 1, CodecProtoData, 0x410001), protoData(ping)...),
			},
			want: &packet.Packet{
				Bean: &Ping{
					Time: "1660360915556",
				},
			},
		},
		{
			name: "json",
			args: args{
				data: append(packet.CreateHeader(jsonHeaderData(ping), 1, CodecJSONData, 0x410001), jsonHeaderData(ping)...),
			},
			want: &packet.Packet{
				Bean: &Ping{
					Time: "1660360915556",
				},
			},
		},
		{
			name: "json no header",
			args: args{
				data: []byte("{\"action\":\"Ping\", \"time\":\"1660360915556\"}"),
			},
			want: &packet.Packet{
				Bean: &Ping{
					Time: "1660360915556",
				},
			},
		},
		{
			name: "not support type",
			args: args{
				data: append(packet.CreateHeader(jsonHeaderData(ping), 1, 8, 0x410001), jsonHeaderData(ping)...),
			},
			wantErr: true,
		},
		{
			name: "not support type",
			args: args{
				data: append(packet.CreateHeader(jsonHeaderData(ping), 1, 5, 0x410001), jsonHeaderData(ping)...),
			},
			wantErr: true,
		},
		{
			name: "not support action key",
			args: args{
				data: append(packet.CreateHeader(jsonHeaderData(ping), 1, CodecJSONData, 0x430001), jsonHeaderData(ping)...),
			},
			wantErr: true,
		},
		{
			name: "not support action key",
			args: args{
				data: []byte("{\"action\":\"Ping2\", \"time\":\"1660360915556\"}"),
			},
			wantErr: true,
		},
		{
			name: "error data",
			args: args{
				data: []byte{0x00, 0x00, 0x00},
			},
			wantErr: true,
		},
		{
			name: "error data",
			args: args{
				data: []byte{0x00},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.Decode(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("CodecHelper.Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got.Bean.(*Ping).Time, tt.want.Bean.(*Ping).Time) {
				t.Errorf("CodecHelper.Decode() = %v, want %v", got.Bean, tt.want.Bean)
			}
		})
	}
}
