package handler

import (
	"fmt"

	"github.com/iotames/easyim/database"
	"github.com/iotames/easyim/model"
)

var HttpApiRoute map[string]func(req *model.Request) error = map[string]func(req *model.Request) error{
	"/api/user/register":      userRegister,
	"/api/user/login":         userLogin,
	"/api/user/logout":        userLogout,
	"/api/user/access_token":  getUserAccessToken,
	"/api/user/refresh_token": getUserRefreshToken,
}

func HttpHandler(req *model.Request) error {
	hreq := req.GetHttpRequest()
	handler, ok := HttpApiRoute[hreq.URL.Path]
	if ok {
		return handler(req)
	}
	data := map[string]interface{}{
		"protocol": hreq.Proto,
		"method":   hreq.Method,
		"host":     hreq.Host,
		"url":      hreq.URL.String(),
		"path":     hreq.URL.Path,
		"query":    hreq.URL.RawQuery,
		"body":     req.GetHttpBody(),
	}
	return model.ResponseApi(data, "找不到路由", 400).Write(*req)
}

type UserLoginForm struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type UserRegisterForm struct {
	UserLoginForm
	Nickname string `json:"nickname"`
}

func userRegister(req *model.Request) error {
	postData := UserRegisterForm{}
	err := req.GetHttpBodyToJson(&postData)
	if err != nil {
		return req.ResponseJson(model.ResponseQueryArgsError(err.Error()))
	}

	u := database.User{Account: postData.Account, Nickname: postData.Nickname}
	u, err = u.Register(postData.Password)
	if err != nil {
		return model.ResponseFail(err.Error(), 400).Write(*req)
	}
	return loginOrRegisterSuccess(u, *req)
}

func userLogin(req *model.Request) error {
	postData := UserLoginForm{}
	err := req.GetHttpBodyToJson(&postData)
	if err != nil {
		return req.ResponseJson(model.ResponseQueryArgsError(err.Error()))
	}
	user := new(database.User)
	user.Account = postData.Account
	database.GetModel(user)
	if user.ID == 0 {
		return model.ResponseFail("找不到登录账号", 400).Write(*req)
	}
	if !user.CheckPassword(postData.Password) {
		return model.ResponseFail("密码不正确", 400).Write(*req)
	}
	return loginOrRegisterSuccess(*user, *req)
}

func loginOrRegisterSuccess(u database.User, req model.Request) error {
	jwtInfo := u.GetJwtInfo()
	data := model.JsonObject{
		"id":            fmt.Sprintf("%d", u.ID),
		"nickname":      u.Nickname,
		"avatar":        database.GetDefaultAvatar(),
		"access_token":  jwtInfo.AccessToken,
		"expires_in":    jwtInfo.Expiresin,
		"refresh_token": jwtInfo.RefreshToken,
	}
	return model.ResponseApi(data, "success", 200).Write(req)
}

func userLogout(req *model.Request) error {
	return nil
}

func getUserAccessToken(req *model.Request) error {
	return nil
}

func getUserRefreshToken(req *model.Request) error {
	return nil
}
