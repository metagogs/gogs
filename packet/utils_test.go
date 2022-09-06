package packet

import (
	"reflect"
	"testing"
)

func TestCreateAction(t *testing.T) {
	type args struct {
		packetType PacketType
		module     uint8
		action     uint16
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{
			name: "Ping",
			args: args{
				packetType: SystemPacket,
				module:     1,
				action:     1,
			},
			want: 0x410001, //0100 0001 - 0000 0000 - 0000 0001
		},
		{
			name: "Pong",
			args: args{
				packetType: SystemPacket,
				module:     1,
				action:     2,
			},
			want: 0x410002, //0100 0001 - 0000 0000 - 0000 0010
		},
		{
			name: "ServicePacket",
			args: args{
				packetType: ServicePacket,
				module:     2,
				action:     2,
			},
			want: 0x820002, //1000 0010 - 0000 0000 - 0000 0010
		},
		{
			name: "ServicePacket2",
			args: args{
				packetType: ServicePacket,
				module:     15,
				action:     3,
			},
			want: 0x8F0003, //1000 1111 - 0000 0000 - 0000 0011
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateAction(tt.args.packetType, tt.args.module, tt.args.action); got != tt.want {
				t.Errorf("CreateAction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestActionToBytes(t *testing.T) {
	type args struct {
		action uint32
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "Ping",
			args: args{
				action: 0x410001,
			},
			want: []byte{0x41, 0x00, 0x01},
		},
		{
			name: "Pong",
			args: args{
				action: 0x410002,
			},
			want: []byte{0x41, 0x00, 0x02},
		},
		{
			name: "ServicePacket",
			args: args{
				action: 0x820002,
			},
			want: []byte{0x82, 0x00, 0x02},
		},
		{
			name: "ServicePacket2",
			args: args{
				action: 0x8F30a3,
			},
			want: []byte{0x8f, 0x30, 0xa3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ActionToBytes(tt.args.action); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ActionToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateHeader(t *testing.T) {
	type args struct {
		data       []byte
		version    uint8
		encodeType uint8
		action     uint32
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "create action",
			args: args{
				data:       []byte{},
				version:    1,
				encodeType: 2,
				action:     0x410002,
			},
			want: []byte{0x7E, 0x0A, 0x41, 0x00, 0x02, 0x00, 0x00, 0x00},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateHeader(tt.args.data, tt.args.version, tt.args.encodeType, tt.args.action); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
