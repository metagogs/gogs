package stringx

import (
	"testing"
)

func BenchmarkStringToBytesPointer(b *testing.B) {
	data := "sfsdfjsdkfjsdklfjskldajflkdsajflkjasdklfjskldjflkdsajflkajsdflkjasdklfj"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StringToBytes(data)
	}
}

func BenchmarkStringToBytes(b *testing.B) {
	data := "sfsdfjsdkfjsdklfjskldajflkdsajflkjasdklfjskldjflkdsajflkajsdflkjasdklfj"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = []byte(data)
	}
}

func BenchmarkBytesToStringPointer(b *testing.B) {
	data := []byte("sfsdfjsdkfjsdklfjskldajflkdsajflkjasdklfjskldjflkdsajflkajsdflkjasdklfj")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BytesToString(data)
	}
}

func BenchmarkBytesToString(b *testing.B) {
	data := []byte("sfsdfjsdkfjsdklfjskldajflkdsajflkjasdklfjskldjflkdsajflkajsdflkjasdklfj")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = string(data)
	}
}
