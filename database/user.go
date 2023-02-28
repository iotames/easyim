package database

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/iotames/miniutils"
)

type JwtInfo struct {
	AccessToken, RefreshToken string
	Expiresin                 int
}

type User struct {
	BaseModel `xorm:"extends"`
	// https://xorm.io/zh/docs/chapter-02/4.columns/  comment	设置字段的注释（当前仅支持mysql）
	Salt         string `xorm:"varchar(64) notnull"`                 //  comment('加密盐')
	PasswordHash string `xorm:"varchar(64) notnull 'password_hash'"` //  comment('密码哈希')
	Account      string `xorm:"varchar(64) notnull unique"`          //  comment('用户名')
	Name         string `xorm:"varchar(32) notnull"`                 //  comment('真实姓名')
	Nickname     string `xorm:"varchar(32) notnull"`                 //  comment('用户昵称')
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

func (u User) Logout(token string) error {
	claims, err := u.CheckTokenGrantType(token, ACCESS_TOKEN)
	if err != nil {
		return err
	}
	u, err = u.GetUserByJwt(token, claims)
	if err != nil {
		return err
	}
	return u.UpdateSalt()
}

func (u *User) UpdateSalt() error {
	u.ResetSalt()
	_, err := UpdateModel(u, map[string]interface{}{"salt": u.Salt})
	return err
}

func (u User) CheckTokenGrantType(token string, match string) (claims map[string]interface{}, err error) {
	claims, err = u.DecodeJwt(token)
	if err != nil {
		return
	}
	grantType, ok := claims[GRANT_TYPE]
	if !ok {
		err = fmt.Errorf("grant_type not found in token claims")
		return
	}
	if grantType.(string) != match {
		err = fmt.Errorf("the token must be " + match)
	}
	return
}

func (u User) getJwt(expiresin int, info map[string]interface{}) string {
	jwt := miniutils.NewJwt(u.Salt)
	token, _ := jwt.Create(info, time.Second*time.Duration(expiresin))
	return token
}

func (u User) DecodeJwt(jwtToken string) (map[string]interface{}, error) {
	jwt := miniutils.NewJwt("")
	return jwt.Decode(jwtToken)
}

// GetUserByJwt 根据JWT字符串获取用户对象。
// claims 可以为 nill
func (u User) GetUserByJwt(jwtStr string, claims map[string]interface{}) (user User, err error) {
	if claims == nil {
		claims, err = u.DecodeJwt(jwtStr)
		if err != nil {
			return
		}
	}
	jsUid := claims["id"].(json.Number)
	var uid int64
	uid, err = jsUid.Int64()
	if err != nil {
		return
	}
	user.ID = uid
	GetModel(&user) // user.Department, user.Position empty
	if user.Account == "" {
		err = fmt.Errorf("user not found")
		user = User{}
		return
	}
	jwt := miniutils.NewJwt(user.Salt)
	_, err = jwt.Parse(jwtStr)
	if err != nil {
		log.Println("--GetUserByJwt--Error:", err)
	}
	return
}

const (
	ACCESS_TOKEN  = "access_token"
	REFRESH_TOKEN = "refresh_token"
	GRANT_TYPE    = "grant_type"
)

func (u User) GetTokenExpiresIn(grantType string) int {
	if grantType == ACCESS_TOKEN {
		return 7200 // 有效期 2 小时
	}
	if grantType == REFRESH_TOKEN {
		return 2592000 // 有效期 30 天. 超长有效期的refresh_token有效防止泄露用户密码
	}
	return 0
}

func (u User) GetAccessToken() string {
	tokenInfo := map[string]interface{}{
		"id":       u.ID,
		"account":  u.Account,
		GRANT_TYPE: ACCESS_TOKEN,
	}
	return u.getJwt(u.GetTokenExpiresIn(ACCESS_TOKEN), tokenInfo)
}

func (u User) GetRefreshToken() string {
	refreshInfo := map[string]interface{}{
		"id":       u.ID,
		"account":  u.Account,
		GRANT_TYPE: REFRESH_TOKEN,
	}
	return u.getJwt(u.GetTokenExpiresIn(REFRESH_TOKEN), refreshInfo)
}

func (u User) GetJwtInfo() JwtInfo {
	return JwtInfo{
		AccessToken:  u.GetAccessToken(),
		RefreshToken: u.GetRefreshToken(),
		Expiresin:    u.GetTokenExpiresIn(ACCESS_TOKEN),
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
		if user.ID > 0 {
			return User{}, fmt.Errorf("登录账号已存在")
		}
	}
	if u.Mobile != "" {
		user.Mobile = u.Mobile
		GetModel(user)
		if user.ID > 0 {
			return User{}, fmt.Errorf("手机号已存在")
		}
	}
	if u.Account == "" && u.Mobile == "" {
		return User{}, fmt.Errorf("登录账号不能为空")
	}
	if u.Account == "" {
		u.Account = u.Mobile
	}
	if u.Mobile == "" {
		u.Mobile = u.Account
	}
	if u.Email == "" {
		u.Email = u.Account + "@example.com"
	}
	if u.Avatar == "" {
		u.Avatar = GetDefaultAvatar()
	}
	user.Account = u.Account
	user.Mobile = u.Mobile
	user.Avatar = u.Avatar
	user.Name = u.Name
	user.Email = u.Email
	user.Nickname = u.Nickname
	user.RemoteAddr = u.RemoteAddr
	user.ResetSalt()
	if password == "" {
		return User{}, fmt.Errorf("注册失败！用户密码不能为空")
	}
	user.SetPasswordHash(password)
	affected, err := CreateModel(user)
	log.Println("affected: ", affected)
	return *user, err
}

func (u *User) ResetSalt() {
	u.Salt = miniutils.GetRandString(64)
}
