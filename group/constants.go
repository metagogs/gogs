package group

import "errors"

var (
	ErrUserNotIntGroup  = errors.New("user not in group")
	ErrUserExistInGroup = errors.New("user exist in group")
)
