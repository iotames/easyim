package handler

import (
	"fmt"
	"time"

	"encoding/json"

	"github.com/iotames/easyim/database"
	"github.com/iotames/easyim/model"
)

var HttpApiRoute map[string]func(req *model.Request) error = map[string]func(req *model.Request) error{
	"/api/user/register":      userRegister,
	"/api/user/login":         userLogin,
	"/api/user/logout":        userLogout,
	"/api/user/check_token":   checkUserToken,
	"/api/user/refresh_token": userRefreshToken,
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

/**
 * @apiDefine LoginParams
 * @apiBody {String} account 用户登录账号
 * @apiBody {String} password 用户登录密码
 */

/**
 * @api {post} /api/user/register 用户注册
 * @apiName 用户注册接口
 * @apiGroup 用户相关
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
 * @apiGroup 用户相关
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
		"avatar":        u.Avatar,
		"access_token":  jwtInfo.AccessToken,
		"expires_in":    jwtInfo.Expiresin,
		"refresh_token": jwtInfo.RefreshToken,
	}
	return model.ResponseApi(data, "success", 200).Write(req)
}

type PostToken struct {
	ResetSecret  bool   `json:"reset_secret"`
	GrantType    string `json:"grant_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

/**
 * @api {post} /api/user/logout 退出登录
 * @apiGroup 用户相关
 * @apiBody {String} access_token 通讯凭证
 * @apiUse PublicCommonParams
 * @apiErrorExample {json} 请求异常示例
 * {"code":400,"msg":"access_token不正确","data":{}}
 * @apiSuccessExample {json} 请求成功示例
 * {"code":200,"msg":"退出登录成功","data":{}}
 */
func userLogout(req *model.Request) error {
	postData := PostToken{}
	err := req.GetHttpBodyToJson(&postData)
	if err != nil {
		return req.ResponseJson(model.ResponseQueryArgsError(err.Error()))
	}
	u := database.User{}
	err = u.Logout(postData.AccessToken)
	if err != nil {
		return model.ResponseFail(err.Error(), 400).Write(*req)
	}
	return model.ResponseOk("登出成功").Write(*req)
}

/**
* @api {get} /api/user/check_token 校验token
* @apiGroup 用户相关
* @apiQuery {String} token 通讯凭证: access_token 或 refresh_token 的值
* @apiSuccess {integer} code 状态码(请求成功为200)
* @apiSuccess {string} msg 请求成功提示信息
* @apiSuccess {Object} data 响应数据
* @apiError {integer} code 请求异常状态码
* @apiError {string} msg 请求异常提示信息
* @apiSuccess {Number} data.expires_in 有效期剩余时间。单位:秒
* @apiSuccess {String} data.grant_type token授权模式。值为 access_token 或 refresh_token
* @apiErrorExample {json} 请求异常示例
* {"code":400,"msg":"access_token不正确","data":{}}
* @apiSuccessExample {json} 请求成功示例
* {"code":200,"msg":"success","data":{"expires_in":7130,"grant_type":"access_token"}}
 */
func checkUserToken(req *model.Request) error {
	query := req.GetHttpRequest().URL.Query()
	token := query.Get("token")
	u := database.User{}
	claims, err := u.DecodeJwt(token)
	if err != nil {
		return model.ResponseFail(err.Error(), 400).Write(*req)
	}
	grantType, ok := claims[database.GRANT_TYPE]
	if !ok {
		return model.ResponseFail("token 不正确", 400).Write(*req)
	}
	_, err = u.GetUserByJwt(token, claims)
	if err != nil {
		return model.ResponseFail(err.Error(), 400).Write(*req)
	}
	expiredAt, _ := claims["exp"].(json.Number).Int64()
	expiresIn := expiredAt - time.Now().Unix()
	data := model.JsonObject{
		"grant_type": grantType,
		"expires_in": expiresIn,
	}
	return model.ResponseApi(data, "success", 200).Write(*req)
}

/**
* @api {post} /api/user/refresh_token 续期token
* @apiGroup 用户相关
* @apiSuccess {integer} code 状态码(请求成功为200)
* @apiSuccess {string} msg 请求成功提示信息
* @apiSuccess {Object} data 响应数据
* @apiError {integer} code 请求异常状态码
* @apiError {string} msg 请求异常提示信息
* @apiBody {String} grant_type 作用域. 值为: access_token或refresh_token. 告诉服务器刷新哪个token.
* @apiBody {String} refresh_token 用于刷新access_token或refresh_token
* @apiBody {Boolean} [reset_secret] 是否重置密钥。默认否false. 重置密钥后，此前使用的access_token和refresh_token都会失效。
* @apiSuccess {Number} data.expires_in 有效期(秒)。从当前时间开始算起。过完失效。
* @apiSuccess {String} data.grant_type 当前token的授权模式。值为 access_token 或 refresh_token
* @apiSuccess {String} data.token 刷新后的 access_token 或 refresh_token 的值
* @apiErrorExample {json} 请求异常示例
* {"code":400,"msg":"token不正确","data":{}}
* @apiSuccessExample {json} 请求成功示例
*{
*    "code": 200,
*    "msg": "success",
*    "data": {
*        "expires_in": 7200,
*        "grant_type": "access_token",
*        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50Ijoid2hxNzgxNjQiLCJleHAiOjE2Nzc1NzYwOTksImdyYW50X3R5cGUiOiJhY2Nlc3NfdG9rZW4iLCJpZCI6MTYzMDM4MTM4ODg5NTA5NjgzMn0.2Tykk-RhXSyy7KoaxqTCEUEOausy3be_QFh-RIY0Lag"
*    }
*}
 */
func userRefreshToken(req *model.Request) error {
	postData := PostToken{}
	err := req.GetHttpBodyToJson(&postData)
	if err != nil {
		return req.ResponseJson(model.ResponseQueryArgsError(err.Error()))
	}
	// 验证refresh_token Begin
	u := new(database.User)
	claims, err := u.CheckTokenGrantType(postData.RefreshToken, database.REFRESH_TOKEN)
	if err != nil {
		return model.ResponseFail(err.Error(), 400).Write(*req)
	}
	uu, err := u.GetUserByJwt(postData.RefreshToken, claims)
	if err != nil {
		return model.ResponseFail(err.Error(), 400).Write(*req)
	}
	// 验证refresh_token End
	if postData.ResetSecret {
		// 重置 salt, 使此前颁发的所有token失效
		u = &uu
		u.UpdateSalt()
	}
	data := model.JsonObject{
		"grant_type": postData.GrantType,
		"expires_in": uu.GetTokenExpiresIn(postData.GrantType),
	}
	if postData.GrantType == database.ACCESS_TOKEN {
		data["token"] = uu.GetAccessToken()
		return model.ResponseApi(data, "success", 200).Write(*req)
	}
	if postData.GrantType == database.REFRESH_TOKEN {
		data["token"] = uu.GetRefreshToken()
		return model.ResponseApi(data, "success", 200).Write(*req)
	}
	return model.ResponseFail("grant_type不正确", 400).Write(*req)
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
