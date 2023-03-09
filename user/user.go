package user

import (
	"strconv"

	"github.com/iotames/easyim/database"
	"github.com/iotames/miniutils"
)

type User struct {
	dbu    database.User
	id     string
	userID int64
}

func NewUser(id string) *User {
	return &User{id: id}
}

func (u *User) parseUserID() int64 {
	uid, err := strconv.Atoi(u.id) // strconv.ParseInt(u.id, 10, 64)
	if err != nil {
		return 0
	}
	u.userID = int64(uid)
	return u.userID
}

func (u *User) getDbu() database.User {
	if u.dbu.ID == 0 {
		uid := u.parseUserID()
		if uid == 0 {
			return u.dbu
		}
		dbu := new(database.User)
		dbu.ID = uid
		database.GetModel(dbu)
		u.dbu = *dbu
	}
	return u.dbu
}

func (u User) GetID() string {
	return u.id
}

func (u *User) CheckToken(token string) bool {
	dbu := u.getDbu()
	if dbu.ID == 0 {
		return false
	}
	jwt := miniutils.NewJwt(dbu.Salt)
	_, err := jwt.Parse(token)
	if err != nil {
		return false
	}
	return true
}
