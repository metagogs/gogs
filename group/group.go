package group

import "context"

type Group interface {
	AddUser(ctx context.Context, uid string) error     // add user to group
	RemoveUser(ctx context.Context, uid string) error  // remove user from group
	RemoveUsers(ctx context.Context, uids []string)    // remove users from group
	RemoveAllUsers(ctx context.Context)                // remove all users from group
	GetUsers(ctx context.Context) []string             // get all users in group
	GetUserCount(ctx context.Context) int              // get user count in group
	GetLastRefresh(ctx context.Context) int64          // get last refresh time
	ContainsUser(ctx context.Context, uid string) bool // check if user is in group
	GetGroupName(ctx context.Context) string           // get group name
	GetGroupID(ctx context.Context) int64              // get group id
}
