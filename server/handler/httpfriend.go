package handler

import (
	"github.com/iotames/easyim/model"
)

/**
 * @api {get} /api/user/friends 好友列表
 * @apiGroup 好友相关
 * @apiDescription 获取用户的好友列表（通讯录显示的联系人）
 * @apiQuery {String} access_token 通讯凭证: access_token 的值
 * @apiSuccess {integer} code 状态码(请求成功为200)
 * @apiSuccess {string} msg 请求成功提示信息
 * @apiSuccess {Object} data 响应数据
 * @apiSuccess {Number} daa.total 好友总数
 * @apiSuccess {Object[]} data.items 用户列表
 * @apiSuccess {String} data.items.id 用户ID
 * @apiSuccess {String} data.items.account 用户账户名
 * @apiSuccess {String} data.items.nickname 用户昵称
 * @apiSuccess {String} data.items.avatar 用户头像
 * @apiError {integer} code 请求异常状态码
 * @apiError {string} msg 请求异常提示信息
 * @apiErrorExample {json} 请求异常示例
 * {"code":400,"msg":"access_token不正确","data":{}}
 * @apiSuccessExample {json} 请求成功示例
 * {"code":200,"msg":"success","data":{}}
 */
func getUserFriends(req *model.Request) error {
	return nil
}

/**
 * @api {post} /api/user/friend/add 发起好友邀请
 * @apiName 发起好友邀请
 * @apiGroup 好友相关
 * @apiDescription 搜索用户后，对该用户发送好友邀请。需要用户同意邀请，才能成为正式好友。
 * @apiBody {String} access_token 发起好友邀请的用户的身份令牌
 * @apiBody {String} to_user_id 被邀请用户的ID
 * @apiUse PublicCommonParams
 * @apiSuccess {String} data.id 本次好友邀请的请求ID
 * @apiErrorExample {json} 请求异常示例
 * {"code":400,"msg":"access_token不正确","data":{}}
 * @apiSuccessExample {json} 请求成功示例
 * {"code":200,"msg":"success","data":{"id":"1629420924912535555"}}
 */
func addUserFriend(req *model.Request) error {
	return nil
}

/**
 * @api {post} /api/user/friend/accept 接受好友邀请
 * @apiGroup 好友相关
 * @apiDescription 同意其他用户发过来的好友申请。正式成为好友。
 * @apiBody {String} access_token 接受好友邀请的用户的身份令牌
 * @apiBody {String} id 好友邀请的请求ID
 * @apiUse PublicCommonParams
 * @apiErrorExample {json} 请求异常示例
 * {"code":400,"msg":"access_token不正确","data":{}}
 * @apiSuccessExample {json} 请求成功示例
 * {"code":200,"msg":"success","data":{}}
 */
func acceptUserFriend(req *model.Request) error {
	return nil
}

/**
 * @api {post} /api/user/friend/remove 删除好友
 * @apiGroup 好友相关
 * @apiDescription TODO
 */
func removeUserFriend(req *model.Request) error {
	return nil
}

/**
 * @api {get} /api/user/search 用户搜索
 * @apiGroup 好友相关
 * @apiDescription 根据用户账号搜索用户 TODO
 */
func searchUser(req *model.Request) error {
	return nil
}
