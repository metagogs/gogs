package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSessionMemory(t *testing.T) {
	m := NewSessionMemory()
	out, exist := m.Get("test")
	assert.False(t, exist)
	assert.Nil(t, out)

	m.Set("test", "test")
	out, exist = m.Get("test")
	assert.True(t, exist)
	assert.Equal(t, "test", out)

	out = m.GetString("test", "")
	assert.Equal(t, "test", out)

	res := m.GetData()
	assert.Equal(t, 1, len(res))
	assert.Equal(t, "test", res["test"])

	m.Delete("test")
	out, exist = m.Get("test")
	assert.False(t, exist)
	assert.Nil(t, out)

	out = m.GetString("test", "default")
	assert.Equal(t, "default", out)

}
