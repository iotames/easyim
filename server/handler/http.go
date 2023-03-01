package handler

import (
	"github.com/iotames/easyim/model"
)

var HttpApiRoute map[string]func(req *model.Request) error = map[string]func(req *model.Request) error{
	"/api/user/register":      userRegister,
	"/api/user/login":         userLogin,
	"/api/user/logout":        userLogout,
	"/api/user/check_token":   checkUserToken,
	"/api/user/refresh_token": userRefreshToken,
	"/api/user/friends":       getUserFriends,
	"/api/user/friend/add":    addUserFriend,
	"/api/user/friend/accept": acceptUserFriend,
	"/api/user/friend/remove": removeUserFriend,
	"/api/user/search":        searchUser,
}

func HttpHandler(req *model.Request) error {
	hreq := req.GetHttpRequest()
	handler, ok := HttpApiRoute[hreq.URL.Path]
	if ok {
		return handler(req)
	}
	data := model.JsonObject{
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

type PostAccessToken struct {
	AccessToken string `json:"access_token"`
}

type PostRefreshToken struct {
	ResetSecret  bool   `json:"reset_secret"`
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

/**
 * @apiDefine PublicCommonParams
 * @apiSuccess {integer} code 状态码(请求成功为200)
 * @apiSuccess {string} msg 请求成功提示信息
 * @apiSuccess {Object} data 响应数据
 * @apiError {integer} code 请求异常状态码
 * @apiError {string} msg 请求异常提示信息
 */

/**
 * @apiDefine ErrorServer
 * @apiErrorExample {json} 请求异常示例
 * {"code":500,"msg":"ServerError.服务器内部错误","data":{}}
 */

/**
 * @apiDefine PublicCommonBlock
 * @apiUse PublicCommonParams
 * @apiUse ErrorServer
 */
