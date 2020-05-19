/*
auth:   wuxun
date:   2020-05-07 18:21
mail:   lbwuxun@qq.com
desc:   how to use or use for what
*/

package _const

var (
	CMD_LOGIN_REQ = 0
	CMD_LOGIN_RSP = 1
	CMD_HEARTBEAT = 2
	CMD_COMMUNICATE = 3
)

//BODY TYPE
var (
	BT_LOGIN_REQ = 1
	BT_LOGIN_RSP = 2
	BT_COMMUNICATE = 3
)
var LOGIN_ACTION = "login"
var LOGIN_FAILED_ACTION = "loginFailed"
var LOGIN_SUCCESS_ACTION = "loginSuccess"
var HEARTBEAT_ACTION = "HeartBeat"

var LOGIN_AUTH = "wuxun"

var MAX_LOSE_HEARTBEAT int = 3
var HEARTBEAT_INTERVAL int = 5

//记录数据长度的大小的字节数量
var LEN_INFO_BYTE_SIZE int = 4
