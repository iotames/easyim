package model

import (
	"github.com/iotames/easyim/contract"
)

type GroupChat struct {
	ID       int64
	msgCount int
	Users    []contract.IUser
}
