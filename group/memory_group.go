package group

import (
	"context"
	"sync"
	"time"

	"github.com/metagogs/gogs/utils/slicex"
)

var _ Group = (*MemoryGroup)(nil)

type MemoryGroup struct {
	mutex       sync.RWMutex
	name        string
	uids        map[string]struct{}
	uidsList    []string
	groupID     int64
	lastRefresh int64
}

func NewMemoryGroup(name string, groupID int64) *MemoryGroup {
	return &MemoryGroup{
		name:        name,
		groupID:     groupID,
		uids:        make(map[string]struct{}),
		lastRefresh: time.Now().Unix(),
	}
}

func (group *MemoryGroup) AddUser(ctx context.Context, uid string) error {
	group.mutex.Lock()
	defer group.mutex.Unlock()
	if _, ok := group.uids[uid]; !ok {
		group.uids[uid] = struct{}{}
		group.uidsList = append(group.uidsList, uid)
		group.lastRefresh = time.Now().Unix()
		return nil
	}

	return ErrUserExistInGroup
}

func (group *MemoryGroup) RemoveUser(ctx context.Context, uid string) error {
	group.mutex.Lock()
	defer group.mutex.Unlock()
	if _, ok := group.uids[uid]; ok {
		delete(group.uids, uid)
		group.uidsList = slicex.RemoveSliceItem(group.uidsList, uid)
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
	group.mutex.Lock()
	defer group.mutex.Unlock()
	group.uids = make(map[string]struct{})
	group.uidsList = nil
}

func (group *MemoryGroup) GetUsers(ctx context.Context) []string {
	group.mutex.RLock()
	defer group.mutex.RUnlock()
	return group.uidsList
}

func (group *MemoryGroup) GetUserCount(ctx context.Context) int {
	return len(group.uidsList)
}

func (group *MemoryGroup) GetLastRefresh(ctx context.Context) int64 {
	return group.lastRefresh
}

func (group *MemoryGroup) ContainsUser(ctx context.Context, uid string) bool {
	group.mutex.RLock()
	defer group.mutex.RUnlock()
	_, ok := group.uids[uid]
	return ok
}

func (group *MemoryGroup) GetGroupName(ctx context.Context) string {
	return group.name
}

func (group *MemoryGroup) GetGroupID(ctx context.Context) int64 {
	return group.groupID
}
