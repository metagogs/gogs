package bytebuffer

import (
	"math/rand"
	"testing"

	"github.com/metagogs/gogs/utils/randstr"
)

//go test -bench=. -benchmem -count=5

func BenchmarkPool(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		generateData, _ := randStr()
		generateData2, _ := randStr()
		bb := defaultPool.Get()
		bb.Write(generateData)
		bb.Write(generateData2)
		defaultPool.Put(bb)
	}
}

func BenchmarkNormal(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		generateData, length := randStr()
		generateData2, length2 := randStr()
		data := make([]byte, length+length2)
		copy(data, generateData)
		copy(data[length:], generateData2)
	}

}

func randStr() ([]byte, int) {
	res := randstr.RandStr(100 + rand.Intn(100))
	data := []byte(res)
	return data, len(data)
}
