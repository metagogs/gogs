package packet

import (
	"reflect"
	"testing"

	"github.com/metagogs/gogs/utils/bytebuffer"
)

func TestParsePacket(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *Packet
		wantErr bool
	}{
		{
			name: "Ping",
			args: args{
				data: []byte{0x7E, 0x0A, 0x41, 0x00, 0x02, 0x00, 0x00, 0x0F, 0x0A, 0x0D, 0x31, 0x36, 0x36, 0x30, 0x33, 0x36, 0x30, 0x39, 0x31, 0x35, 0x35, 0x35, 0x36},
			},
			wantErr: false,
			want: &Packet{
				HeaderByte: []byte{0x7E, 0x0A, 0x41, 0x00, 0x02, 0x00, 0x00, 0x0F},
				Data:       []byte{0x0A, 0x0D, 0x31, 0x36, 0x36, 0x30, 0x33, 0x36, 0x30, 0x39, 0x31, 0x35, 0x35, 0x35, 0x36},
			},
		},
		{
			name: "PingErr",
			args: args{
				data: []byte{0x8E, 0x0A, 0x41, 0x00, 0x02, 0x00, 0x00, 0x0F, 0x0A, 0x0D, 0x31, 0x36, 0x36, 0x30, 0x33, 0x36, 0x30, 0x39, 0x31, 0x35, 0x35, 0x35, 0x36},
			},
			wantErr: true,
		},
		{
			name: "PingErr",
			args: args{
				data: []byte{0x7E, 0x0A, 0x11, 0x00, 0x02, 0x00, 0x00, 0x0F, 0x0A, 0x0D, 0x31, 0x36, 0x36, 0x30, 0x33, 0x36, 0x30, 0x39, 0x31, 0x35, 0x35, 0x35, 0x36},
			},
			wantErr: true,
		},
		{
			name: "Ping",
			args: args{
				data: []byte{0x7E},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePacket(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePacket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParsePacket() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_GetActionKey(t *testing.T) {
	type fields struct {
		HeaderByte []byte
		Data       []byte
		Bean       interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{
			name: "get action key",
			fields: fields{
				HeaderByte: []byte{0x7E, 0x0A, 0x41, 0x00, 0x02, 0x00, 0x00, 0x0F},
			},
			want: 0x410002,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				HeaderByte: tt.fields.HeaderByte,
				Data:       tt.fields.Data,
				Bean:       tt.fields.Bean,
			}
			if got := p.GetActionKey(); got != tt.want {
				t.Errorf("Packet.GetActionKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_GetPacketType(t *testing.T) {
	type fields struct {
		HeaderByte []byte
		Data       []byte
		Bean       interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   PacketType
	}{
		{
			name: "get packet type",
			fields: fields{
				HeaderByte: []byte{0x7E, 0x0A, 0x41, 0x00, 0x02, 0x00, 0x00, 0x0F},
			},
			want: SystemPacket,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				HeaderByte: tt.fields.HeaderByte,
				Data:       tt.fields.Data,
				Bean:       tt.fields.Bean,
			}
			if got := p.GetPacketType(); got != tt.want {
				t.Errorf("Packet.GetPacketType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_GetVersion(t *testing.T) {
	type fields struct {
		HeaderByte []byte
		Data       []byte
		Bean       interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   uint8
	}{
		{
			name: "get packet version",
			fields: fields{
				HeaderByte: []byte{0x7E, 0x0A, 0x41, 0x00, 0x02, 0x00, 0x00, 0x0F},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				HeaderByte: tt.fields.HeaderByte,
				Data:       tt.fields.Data,
				Bean:       tt.fields.Bean,
			}
			if got := p.GetVersion(); got != tt.want {
				t.Errorf("Packet.GetVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_GetEncodeType(t *testing.T) {
	type fields struct {
		HeaderByte []byte
		Data       []byte
		Bean       interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   uint8
	}{
		{
			name: "get packet version",
			fields: fields{
				HeaderByte: []byte{0x7E, 0x0A, 0x41, 0x00, 0x02, 0x00, 0x00, 0x0F},
			},
			want: uint8(2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				HeaderByte: tt.fields.HeaderByte,
				Data:       tt.fields.Data,
				Bean:       tt.fields.Bean,
			}
			if got := p.GetEncodeType(); got != tt.want {
				t.Errorf("Packet.GetEncodeType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_GetModule(t *testing.T) {
	type fields struct {
		HeaderByte []byte
		Data       []byte
		Bean       interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   uint8
	}{
		{
			name: "get module",
			fields: fields{
				HeaderByte: []byte{0x7E, 0x0A, 0x41, 0x00, 0x02, 0x00, 0x00, 0x0F},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				HeaderByte: tt.fields.HeaderByte,
				Data:       tt.fields.Data,
				Bean:       tt.fields.Bean,
			}
			if got := p.GetModule(); got != tt.want {
				t.Errorf("Packet.GetModule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_GetLength(t *testing.T) {
	type fields struct {
		HeaderByte []byte
		Data       []byte
		Bean       interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{
			name: "get length",
			fields: fields{
				HeaderByte: []byte{0x7E, 0x0A, 0x41, 0x00, 0x02, 0x00, 0x00, 0x0F},
			},
			want: 0x0f,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				HeaderByte: tt.fields.HeaderByte,
				Data:       tt.fields.Data,
				Bean:       tt.fields.Bean,
			}
			if got := p.GetLength(); got != tt.want {
				t.Errorf("Packet.GetLength() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildPacketWithHeader(t *testing.T) {
	type args struct {
		data       []byte
		version    uint8
		ecnodeType uint8
		action     uint32
	}
	tests := []struct {
		name string
		args args
		want *Packet
	}{
		{
			name: "build",
			args: args{
				data:       []byte{0x0A, 0x0D, 0x31, 0x36, 0x36, 0x30, 0x33, 0x36, 0x30, 0x39, 0x31, 0x35, 0x35, 0x35, 0x36},
				version:    1,
				ecnodeType: 2,
				action:     0x410002,
			},
			want: &Packet{
				HeaderByte: []byte{0x7E, 0x0A, 0x41, 0x00, 0x02, 0x00, 0x00, 0x0F},
				Data:       []byte{0x0A, 0x0D, 0x31, 0x36, 0x36, 0x30, 0x33, 0x36, 0x30, 0x39, 0x31, 0x35, 0x35, 0x35, 0x36},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPacketWithHeader(tt.args.data, tt.args.version, tt.args.ecnodeType, tt.args.action); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildPacketWithHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPacket_ToData(t *testing.T) {
	type fields struct {
		HeaderByte []byte
		Data       []byte
		Bean       interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   *bytebuffer.ByteBuffer
	}{
		{
			name: "to data",
			fields: fields{
				HeaderByte: []byte{0x7E, 0x0A, 0x41, 0x00, 0x02, 0x00, 0x00, 0x0F},
				Data:       []byte{0x0A, 0x0D, 0x31, 0x36, 0x36, 0x30, 0x33, 0x36, 0x30, 0x39, 0x31, 0x35, 0x35, 0x35, 0x36},
			},
			want: &bytebuffer.ByteBuffer{
				B: []byte{0x7E, 0x0A, 0x41, 0x00, 0x02, 0x00, 0x00, 0x0F, 0x0A, 0x0D, 0x31, 0x36, 0x36, 0x30, 0x33, 0x36, 0x30, 0x39, 0x31, 0x35, 0x35, 0x35, 0x36},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{
				HeaderByte: tt.fields.HeaderByte,
				Data:       tt.fields.Data,
				Bean:       tt.fields.Bean,
			}
			if got := p.ToData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Packet.ToData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPacket(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want *Packet
	}{
		{
			name: "to data",
			args: args{
				data: []byte{0x0A, 0x0D, 0x31, 0x36, 0x36, 0x30, 0x33, 0x36, 0x30, 0x39, 0x31, 0x35, 0x35, 0x35, 0x36},
			},
			want: &Packet{
				HeaderByte: []byte{},
				Data:       []byte{0x0A, 0x0D, 0x31, 0x36, 0x36, 0x30, 0x33, 0x36, 0x30, 0x39, 0x31, 0x35, 0x35, 0x35, 0x36},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPacket(tt.args.data); !reflect.DeepEqual(got.Data, tt.want.Data) {
				t.Errorf("NewPacket() = %v, want %v", got, tt.want)
			}
		})
	}
}
