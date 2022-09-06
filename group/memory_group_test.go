package group

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryGroup_AddUser(t *testing.T) {
	memGroup := NewMemoryGroup("name", 1)
	id := memGroup.GetGroupID(nil)
	assert.Equal(t, int64(1), id)
	name := memGroup.GetGroupName(nil)
	assert.Equal(t, "name", name)

	err := memGroup.AddUser(nil, "test")
	assert.Nil(t, err)

	err = memGroup.AddUser(nil, "test")
	assert.Equal(t, ErrUserExistInGroup, err)

	err = memGroup.AddUser(nil, "test2")
	assert.Nil(t, err)

	uids := memGroup.GetUsers(nil)
	assert.Equal(t, 2, len(uids))
	assert.Contains(t, uids, "test")
	assert.Contains(t, uids, "test2")

	count := memGroup.GetUserCount(nil)
	assert.Equal(t, 2, count)

	err = memGroup.RemoveUser(nil, "test")
	assert.Nil(t, err)

	uids = memGroup.GetUsers(nil)
	assert.Equal(t, 1, len(uids))
	assert.Contains(t, uids, "test2")

	memGroup.RemoveAllUsers(nil)
	uids = memGroup.GetUsers(nil)
	assert.Equal(t, 0, len(uids))

	count = memGroup.GetUserCount(nil)
	assert.Equal(t, 0, count)

	err = memGroup.RemoveUser(nil, "test")
	assert.Equal(t, ErrUserNotIntGroup, err)

	err = memGroup.AddUser(nil, "test")
	assert.Nil(t, err)

	err = memGroup.AddUser(nil, "test2")
	assert.Nil(t, err)

	memGroup.RemoveUsers(nil, []string{"test"})
	uids = memGroup.GetUsers(nil)
	assert.Equal(t, 1, len(uids))
	assert.Contains(t, uids, "test2")

	exist := memGroup.ContainsUser(nil, "test")
	assert.False(t, exist)

	exist = memGroup.ContainsUser(nil, "test2")
	assert.True(t, exist)

}
