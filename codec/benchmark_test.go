package codec

import (
	"encoding/json"
	"testing"

	"github.com/metagogs/gogs/config"
)

// go test -bench=. -benchmem -benchtime=20s
// go test -bench=BenchmarkEncodePacketProto -run=none -benchmem
// go test -bench=BenchmarkEncodePacketProto -run=none -benchmem -memprofile=mem.pprof -cpuprofile=cpu.pprof
// go tool pprof -http=:8081 cpu.pprof
func BenchmarkEncodePacketProto(b *testing.B) {
	dispatch := TestDispatch{}
	d := NewCodecHelper(config.NewDefaultConfig(), &dispatch)
	d.RegisterEncode(CodecProtoData, &ProtoEncode{})
	pingTest1 := &Ping{
		Time: "123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = d.Encode(pingTest1)
	}
}

func BenchmarkEncodePacketJson(b *testing.B) {
	dispatch := TestDispatch{}
	d := NewCodecHelper(config.NewDefaultConfig(), &dispatch)
	d.RegisterEncode(CodecJSONData, &JSONEncode{})
	pingTest1 := &Ping{
		Time: "123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = d.Encode(pingTest1)
	}
}

func BenchmarkEncodePacketJsonNoHeader(b *testing.B) {
	dispatch := TestDispatch{}
	d := NewCodecHelper(config.NewDefaultConfig(), &dispatch)
	d.RegisterEncode(CodecJSONDataNoHeader, &JSONEncode{})
	pingTest1 := &Ping{
		Time: "123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = d.Encode(pingTest1)
	}
}

func BenchmarkJSONOriginEncode(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(&Ping{
			Time: "1234567890",
		})
	}
}

func BenchmarkJSONEncode(b *testing.B) {
	encode := &JSONEncode{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = encode.Encode(&Ping{
			Time: "1234567890",
		})
	}
}

func BenchmarkProtoEncode(b *testing.B) {
	encode := &ProtoEncode{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = encode.Encode(&Ping{
			Time: "1234567890",
		})
	}
}

func BenchmarkJSONOriginDecode(b *testing.B) {
	data, _ := json.Marshal(&Ping{
		Time: "1234567890",
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = json.Unmarshal(data, &Ping{})
	}
}

func BenchmarkJSONDecode(b *testing.B) {
	encode := &JSONEncode{}
	decode := &JSONDecode{}

	data, _ := encode.Encode(&Ping{
		Time: "1234567890",
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = decode.Decode(data, &Ping{})
	}
}

func BenchmarkProtoDecode(b *testing.B) {
	encode := &ProtoEncode{}
	decode := &ProtoDecode{}

	data, _ := encode.Encode(&Ping{
		Time: "1234567890",
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = decode.Decode(data, &Ping{})
	}
}
