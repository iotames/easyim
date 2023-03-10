package handler

import (
	"fmt"
	"time"

	"encoding/json"

	"github.com/iotames/easyim/database"
	"github.com/iotames/easyim/model"
)

/**
 * @apiDefine LoginParams
 * @apiBody {String} account 用户登录账号
 * @apiBody {String} password 用户登录密码
 */

/**
 * @api {post} /api/user/login 用户登录
 * @apiName 用户登录
 * @apiGroup 用户与权限
 * @apiUse LoginParams
 * @apiUse PublicCommonParams
 * @apiErrorExample {json} 请求异常示例
 * {"code":400,"msg":"密码不正确","data":{}}
 * @apiUse LoginOrRegisterSuccessBlock
 */
func userLogin(req *model.Request, resp *model.Response) model.Response {
	postData := UserLoginForm{}
	err := req.GetHttpBodyToJson(&postData)
	if err != nil {
		return resp.Json(model.ResponseQueryArgsError(err.Error()))
	}
	user := new(database.User)
	user.Account = postData.Account
	database.GetModel(user)
	if user.ID == 0 {
		return resp.Json(model.ResponseFail("找不到登录账号", 400))
	}
	if !user.CheckPassword(postData.Password) {
		return resp.Json(model.ResponseFail("密码不正确", 400))
	}
	return loginOrRegisterSuccess(*user, resp)
}

/**
 * @api {post} /api/user/register 用户注册
 * @apiName 用户注册接口
 * @apiGroup 用户与权限
 * @apiUse LoginParams
 * @apiBody {String} nickname 用户昵称
 * @apiUse PublicCommonParams
 * @apiErrorExample {json} 请求异常示例
 * {"code":400,"msg":"注册失败！登录账号已存在","data":{}}
 * @apiUse LoginOrRegisterSuccessBlock
 */
func userRegister(req *model.Request, resp *model.Response) model.Response {
	postData := UserRegisterForm{}
	err := req.GetHttpBodyToJson(&postData)
	if err != nil {
		return resp.Json(model.ResponseQueryArgsError(err.Error()))
	}
	u := database.User{Account: postData.Account, Nickname: postData.Nickname}
	u, err = u.Register(postData.Password)
	if err != nil {
		return resp.Json(model.ResponseFail(err.Error(), 400))
	}
	return loginOrRegisterSuccess(u, resp)
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
func loginOrRegisterSuccess(u database.User, resp *model.Response) model.Response {
	jwtInfo := u.GetJwtInfo()
	data := model.JsonObject{
		"id":            fmt.Sprintf("%d", u.ID),
		"nickname":      u.Nickname,
		"avatar":        u.Avatar,
		"access_token":  jwtInfo.AccessToken,
		"expires_in":    jwtInfo.Expiresin,
		"refresh_token": jwtInfo.RefreshToken,
	}
	return resp.Json(model.ResponseApi(data, "success", 200))
}

/**
 * @api {post} /api/user/logout 退出登录
 * @apiGroup 用户与权限
 * @apiBody {String} access_token 通讯凭证
 * @apiUse PublicCommonParams
 * @apiErrorExample {json} 请求异常示例
 * {"code":400,"msg":"access_token不正确","data":{}}
 * @apiSuccessExample {json} 请求成功示例
 * {"code":200,"msg":"退出登录成功","data":{}}
 */
func userLogout(req *model.Request, resp *model.Response) model.Response {
	postData := PostAccessToken{}
	err := req.GetHttpBodyToJson(&postData)
	if err != nil {
		return resp.Json(model.ResponseQueryArgsError(err.Error()))
	}
	u := database.User{}
	err = u.Logout(postData.AccessToken)
	if err != nil {
		return resp.Json(model.ResponseFail(err.Error(), 400))
	}
	return resp.Json(model.ResponseOk("登出成功"))
}

/**
 * @api {get} /api/user/check_token 校验token
 * @apiGroup 用户与权限
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
func checkUserToken(req *model.Request, resp *model.Response) model.Response {
	query := req.GetHttpRequest().URL.Query()
	token := query.Get("token")
	u := database.User{}
	claims, err := u.DecodeJwt(token)
	if err != nil {
		return resp.Json(model.ResponseFail(err.Error(), 400))
	}
	grantType, ok := claims[database.GRANT_TYPE]
	if !ok {
		return resp.Json(model.ResponseFail("token 不正确", 400))
	}
	_, err = u.GetUserByJwt(token, claims)
	if err != nil {
		return resp.Json(model.ResponseFail(err.Error(), 400))
	}
	expiredAt, _ := claims["exp"].(json.Number).Int64()
	expiresIn := expiredAt - time.Now().Unix()
	data := model.JsonObject{
		"grant_type": grantType,
		"expires_in": expiresIn,
	}
	return resp.Json(model.ResponseApi(data, "success", 200))
}

/**
* @api {post} /api/user/refresh_token 续期token
* @apiGroup 用户与权限
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
func userRefreshToken(req *model.Request, resp *model.Response) model.Response {
	postData := PostRefreshToken{}
	err := req.GetHttpBodyToJson(&postData)
	if err != nil {
		return resp.Json(model.ResponseQueryArgsError(err.Error()))
	}
	// 验证refresh_token Begin
	u := new(database.User)
	claims, err := u.CheckTokenGrantType(postData.RefreshToken, database.REFRESH_TOKEN)
	if err != nil {
		return resp.Json(model.ResponseFail(err.Error(), 400))
	}
	uu, err := u.GetUserByJwt(postData.RefreshToken, claims)
	if err != nil {
		return resp.Json(model.ResponseFail(err.Error(), 400))
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
		return resp.Json(model.ResponseApi(data, "success", 200))
	}
	if postData.GrantType == database.REFRESH_TOKEN {
		data["token"] = uu.GetRefreshToken()
		return resp.Json(model.ResponseApi(data, "success", 200))
	}
	return resp.Json(model.ResponseFail("grant_type不正确", 400))
}
