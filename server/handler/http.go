package handler

import (
	"fmt"

	"github.com/iotames/easyim/model"
)

var HttpApiRoute map[string]func(req *model.Request, resp *model.Response) model.Response = map[string]func(req *model.Request, resp *model.Response) model.Response{
	"POST /api/user/register":      userRegister,
	"POST /api/user/login":         userLogin,
	"POST /api/user/logout":        userLogout,
	"GET /api/user/check_token":    checkUserToken,
	"POST /api/user/refresh_token": userRefreshToken,
	"GET /api/user/friends":        getUserFriends,
	"POST /api/user/friend/add":    addUserFriend,
	"POST /api/user/friend/accept": acceptUserFriend,
	"POST /api/user/friend/remove": removeUserFriend,
	"GET /api/user/search":         searchUser,
	"POST /api/local/stop":         closeListener,
}

func HttpHandler(req *model.Request) error {
	hreq := req.GetHttpRequest()
	resp := model.NewResponse(req.GetConn())
	if hreq.Method == "OPTIONS" {
		return resp.OPTIONS().Write()
	}
	handler, ok := HttpApiRoute[fmt.Sprintf("%s %s", hreq.Method, hreq.URL.Path)]
	if ok {
		resp := handler(req, resp)
		return resp.Write()
	}
	return HttpNotFound(req, resp).Write()
}

func HttpNotFound(req *model.Request, resp *model.Response) model.Response {
	hreq := req.GetHttpRequest()
	data := model.JsonObject{
		"protocol": hreq.Proto,
		"method":   hreq.Method,
		"host":     hreq.Host,
		"url":      hreq.URL.String(),
		"path":     hreq.URL.Path,
		"query":    hreq.URL.RawQuery,
		"body":     req.GetHttpBody(),
	}
	return resp.Json(model.ResponseApi(data, "找不到路由", 400))
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
