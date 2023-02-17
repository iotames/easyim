package database

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/iotames/miniutils"
)

type JwtInfo struct {
	Token     string
	Expiresin int
}

type User struct {
	BaseModel `xorm:"extends"`
	// https://xorm.io/zh/docs/chapter-02/4.columns/  comment	设置字段的注释（当前仅支持mysql）
	Salt         string `xorm:"varchar(64) notnull"`                 //  comment('加密盐')
	PasswordHash string `xorm:"varchar(64) notnull 'password_hash'"` //  comment('密码哈希')
	Account      string `xorm:"varchar(64) notnull unique"`          //  comment('用户名')
	Name         string `xorm:"varchar(32) notnull"`                 //  comment('真实姓名')
	Mobile       string `xorm:"varchar(32) notnull unique"`          // comment('手机号')
	Email        string `xorm:"varchar(32) notnull unique"`          //  comment('电子邮箱')
	Avatar       string `xorm:"varchar(500) notnull"`                // comment('用户头像')
	RemoteAddr   string `xorm:"varchar(32) notnull 'remote_addr'"`   //  comment('客户端地址')
}

func (u User) TableName() string {
	return "users"
}

func GetDefaultAvatar() string {
	return "https://gw.alipayobjects.com/zos/antfincdn/XAosXuNZyF/BiazfanxmamNRoxxVxka.png"
}

func (u User) getJwt(expiresin int) string {
	jwt := miniutils.NewJwt(u.Salt)
	info := map[string]interface{}{
		"id":      u.ID,
		"account": u.Account,
		"avatar":  u.Avatar,
	}
	token, _ := jwt.Create(info, time.Second*time.Duration(expiresin))
	return token
}

func (u User) GetUserByJwt(jwtStr string) (user User, err error) {
	var segInfo map[string]interface{}
	jwt := miniutils.NewJwt("")
	segInfo, err = jwt.Decode(jwtStr)
	if err != nil {
		return
	}
	jsUid := segInfo["id"].(json.Number)
	uid, _ := jsUid.Int64()
	user.ID = uid
	GetModel(&user) // user.Department, user.Position empty
	log.Println("---FoundUser--By--Jwt---user.Salt------", user.Salt)
	jwt = miniutils.NewJwt(user.Salt)
	_, err = jwt.Parse(jwtStr)
	if err != nil {
		log.Println("--GetUserByJwt--Error:", err)
	}
	return
}

func (u User) GetJwtInfo() JwtInfo {
	expiresin := 3600 * 24 * 7 // 有效期 7 天
	return JwtInfo{
		Token:     u.getJwt(expiresin),
		Expiresin: expiresin,
	}
}

func (u User) getPasswordHash(password string) string {
	return miniutils.GetSha256(miniutils.GetSha256(password))
}

func (u *User) SetPasswordHash(password string) {
	u.PasswordHash = u.getPasswordHash(password)
}

func (u User) CheckPassword(password string) bool {
	return u.PasswordHash == u.getPasswordHash(password)
}

func (u User) Register(password string) (User, error) {
	user := new(User)
	if u.Account != "" {
		user.Account = u.Account
		GetModel(user)
	}
	if u.Mobile != "" {
		user.Mobile = u.Mobile
		GetModel(user)
	}
	if user.ID > 0 {
		return User{}, fmt.Errorf("error: Regiser Fail. User exists")
	}
	user.Account = u.Account
	user.Mobile = u.Mobile
	user.Avatar = u.Avatar
	user.Name = u.Name
	user.RemoteAddr = u.RemoteAddr
	user.ResetSalt()
	if password == "" {
		return User{}, fmt.Errorf("error: Regiser Fail. User password can not be empty")
	}
	user.SetPasswordHash(password)
	affected, err := CreateModel(user)
	log.Println("affected: ", affected)
	return *user, err
}

func (u *User) ResetSalt() {
	u.Salt = miniutils.GetRandString(64)
}
