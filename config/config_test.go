package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	os.Setenv("GOGS_GOPPROFADDR", "9996")

	config := NewConfig()
	assert.Equal(t, true, config.Debug)
	assert.Equal(t, 9996, config.GopprofAddr)
	assert.Equal(t, 19, config.AgentHeartBeatTimeout)

	config = NewConfig("config2.yaml")
	assert.Equal(t, true, config.Debug)
	assert.Equal(t, 9996, config.GopprofAddr)
	assert.Equal(t, 19, config.AgentHeartBeatTimeout)
}
