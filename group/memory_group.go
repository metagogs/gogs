package group

import (
	"context"
	"sync"
	"time"
)

var _ Group = (*MemoryGroup)(nil)

type MemoryGroup struct {
	name        string
	uids        sync.Map
	groupID     int64
	lastRefresh int64
}

func NewMemoryGroup(name string, groupID int64) *MemoryGroup {
	return &MemoryGroup{
		name:        name,
		groupID:     groupID,
		lastRefresh: time.Now().Unix(),
	}
}

func (group *MemoryGroup) AddUser(ctx context.Context, uid string) error {
	if _, ok := group.uids.Load(uid); !ok {
		group.uids.Store(uid, uid)
		group.lastRefresh = time.Now().Unix()
		return nil
	}

	return ErrUserExistInGroup
}

func (group *MemoryGroup) RemoveUser(ctx context.Context, uid string) error {
	if _, ok := group.uids.Load(uid); ok {
		group.uids.Delete(uid)
		group.lastRefresh = time.Now().Unix()
		return nil
	}

	return ErrUserNotIntGroup
}

func (group *MemoryGroup) RemoveUsers(ctx context.Context, uids []string) {
	for _, uid := range uids {
		_ = group.RemoveUser(ctx, uid)
		group.lastRefresh = time.Now().Unix()
	}
}

func (group *MemoryGroup) RemoveAllUsers(ctx context.Context) {
	group.uids.Range(func(key, value interface{}) bool {
		group.uids.Delete(key)
		group.lastRefresh = time.Now().Unix()
		return true
	})
}

func (group *MemoryGroup) GetUsers(ctx context.Context) []string {
	uids := []string{}
	group.uids.Range(func(key, value interface{}) bool {
		uids = append(uids, key.(string))
		return true
	})

	return uids
}

func (group *MemoryGroup) GetUserCount(ctx context.Context) int {
	count := 0
	group.uids.Range(func(key, value interface{}) bool {
		count++
		return true
	})

	return count
}

func (group *MemoryGroup) GetLastRefresh(ctx context.Context) int64 {
	return group.lastRefresh
}

func (group *MemoryGroup) ContainsUser(ctx context.Context, uid string) bool {
	_, ok := group.uids.Load(uid)
	return ok
}

func (group *MemoryGroup) GetGroupName(ctx context.Context) string {
	return group.name
}

func (group *MemoryGroup) GetGroupID(ctx context.Context) int64 {
	return group.groupID
}
