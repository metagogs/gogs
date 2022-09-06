package bytebuffer

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPool_Get(t *testing.T) {
	got := Get()
	assert.NotNil(t, got)
	got.Write([]byte("hello"))
	assert.Equal(t, "hello", got.String())
	got.WriteString(" world")
	assert.Equal(t, "hello world", got.String())
	got.Set([]byte("foo"))
	assert.Equal(t, "foo", got.String())
	assert.Equal(t, 3, got.Len())
	got.SetString("bar")
	assert.Equal(t, "bar", got.String())
	got.WriteByte('!')
	assert.Equal(t, "bar!", got.String())
	Put(got)
	got2 := Get()
	got2.ReadFrom(strings.NewReader("hello"))
	assert.Equal(t, "hello", got2.String())
	assert.Equal(t, []byte("hello"), got2.Bytes())
	Put(got2)
	assert.Equal(t, 0, got.Len())
	assert.Equal(t, 0, got2.Len())
}
