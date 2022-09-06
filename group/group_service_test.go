package group

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroupServer_CreateMemoryGroup(t *testing.T) {
	groupServer := NewGroupServer()

	groupServer.CreateMemoryGroup("test")
	groupServer.CreateMemoryGroup("test")
	g, exist := groupServer.GetGroup("test")
	assert.True(t, exist)
	assert.NotNil(t, g)

	g, exist = groupServer.GetGroup("test2")
	assert.False(t, exist)
	assert.Nil(t, g)

	groupServer.DeletGroup("test")
	g, exist = groupServer.GetGroup("test")
	assert.False(t, exist)
	assert.Nil(t, g)

	groupServer.CreateMemoryGroup("test")
	g, exist = groupServer.GetGroup("test")
	assert.True(t, exist)
	assert.NotNil(t, g)

	err := g.AddUser(nil, "user1")
	assert.Nil(t, err)

	err = g.AddUser(nil, "user2")
	assert.Nil(t, err)

	groupServer.CreateMemoryGroup("test2")
	g2, exist := groupServer.GetGroup("test2")
	assert.True(t, exist)
	assert.NotNil(t, g2)

	err = g2.AddUser(nil, "user1")
	assert.Nil(t, err)

	groupServer.DeleteUserByName("user1")
	exist = g.ContainsUser(nil, "user1")
	assert.False(t, exist)

	exist = g2.ContainsUser(nil, "user1")
	assert.False(t, exist)

	exist = g.ContainsUser(nil, "user2")
	assert.True(t, exist)

	exist = g2.ContainsUser(nil, "user2")
	assert.False(t, exist)
}
