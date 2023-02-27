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

/**
 * @apiDefine LoginParams
 * @apiBody {String} account 用户登录账号
 * @apiBody {String} password 用户登录密码
 */

/**
 * @api {post} /api/user/register 用户注册
 * @apiName 用户注册接口
 * @apiGroup 注册登录
 * @apiUse LoginParams
 * @apiBody {String} nickname 用户昵称
 * @apiUse PublicCommonParams
 * @apiErrorExample {json} 请求异常示例
 * {"code":400,"msg":"注册失败！登录账号已存在","data":{}}
 * @apiUse LoginOrRegisterSuccessBlock
 */
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

/**
 * @api {post} /api/user/login 用户登录
 * @apiName 用户登录
 * @apiGroup 注册登录
 * @apiUse LoginParams
 * @apiUse PublicCommonParams
 * @apiErrorExample {json} 请求异常示例
 * {"code":400,"msg":"密码不正确","data":{}}
 * @apiUse LoginOrRegisterSuccessBlock
 */
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

/**
 * @apiDefine LoginOrRegisterSuccessBlock
 * @apiSuccess {String} data.id 用户ID
 * @apiSuccess {String} data.nickname 用户昵称
 * @apiSuccess {String} data.avatar 用户头像
 * @apiSuccess {String} data.access_token 通讯凭证，在IM通讯中或需要身份鉴权的API中使用。 有效期7200秒(2小时)
 * @apiSuccess {String} data.expires_in access_token(通讯凭证)的有效期(一般为7200秒)
 * @apiSuccess {String} data.refresh_token 用于刷新access_token(通讯凭证)，有效期30天
 * @apiSuccessExample {json} 请求成功示例
 * {"code":200,"msg":"success","data":{"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50IjoiYWNjb3VudFd1aGFucWluZyIsImF2YXRhciI6IiIsImV4cCI6MTY3NzQ3NDkxNCwiaWQiOjE2Mjk0MjA5MjQ5MTI1Mzc2MDB9.IucceY2x7FSB81-nxEj_yMYggYaBnCzEX1GA8LdzPCE","avatar":"https://gw.alipayobjects.com/zos/antfincdn/XAosXuNZyF/BiazfanxmamNRoxxVxka.png","expires_in":7200,"id":"1629420924912537600","nickname":"飞天的猪","refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50IjoiYWNjb3VudFd1aGFucWluZyIsImF2YXRhciI6IiIsImV4cCI6MTY4MDA1OTcxNCwiaWQiOjE2Mjk0MjA5MjQ5MTI1Mzc2MDB9.Mu4N8-uegCdq26ocIx7HINmoUgLyrwpqo4cYDslHwzs"}}
 */
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

type UserLogoutData struct {
	Account     string `json:"account"`
	AccessToken string `json:"access_token"`
}

/**
 * @api {post} /api/user/logout 退出登录
 * @apiGroup 注册登录
 * @apiBody {String} access_token 通讯凭证
 * @apiUse PublicCommonParams
 * @apiErrorExample {json} 请求异常示例
 * {"code":400,"msg":"access_token不正确","data":{}}
 * @apiSuccessExample {json} 请求成功示例
 * {}
 */
func userLogout(req *model.Request) error {
	postData := UserLogoutData{}
	err := req.GetHttpBodyToJson(&postData)
	if err != nil {
		return req.ResponseJson(model.ResponseQueryArgsError(err.Error()))
	}

	u := database.User{}
	u, err = u.GetUserByJwt(postData.AccessToken)
	if err != nil {
		return model.ResponseFail(err.Error(), 400).Write(*req)
	}
	u.ResetSalt()
	_, err = database.UpdateModel(&u, map[string]interface{}{"salt": u.Salt})
	if err != nil {
		return model.ResponseServerError().Write(*req)
	}
	return model.ResponseOk("退出登录成功").Write(*req)
}

func getUserAccessToken(req *model.Request) error {
	return nil
}

func getUserRefreshToken(req *model.Request) error {
	return nil
}

/**
 * @apiDefine PublicCommonParams
 * @apiSuccess {integer} code 状态码(请求成功为200)
 * @apiSuccess {string} msg 请求成功提示信息
 * @apiSuccess {Object} data 响应数据
 * @apiSuccess {integer} code 请求异常状态码
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
