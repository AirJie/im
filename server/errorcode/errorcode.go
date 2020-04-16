package errorcode

import (
	"errors"
)

var (
	ErrUnknown = errors.New("error unknown ")
	ErrUnauthorized = errors.New("error unauthorized")
	ErrNotInGroup = errors.New("error not in group")
	ErrDeviceNotBindUser = errors.New("error device not bind user")
	ErrBadRequest = errors.New("error bad request")
	ErrUserAlreadyExist = errors.New("error user already exist")
	ErrGroupAlreadyExist = errors.New("error group already exist")
	ErrUserNotExist = errors.New("error user not exists")
)

