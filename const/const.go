/*
auth:   wuxun
date:   2020-05-07 18:21
mail:   lbwuxun@qq.com
desc:   how to use or use for what
*/

package _const

// 心跳丢失的最大次数，超过此次数就算掉线
var MAX_LOSE_HEARTBEAT int = 3
// 心跳的间隔秒数
var HEARTBEAT_INTERVAL int = 5

//记录数据长度的大小的字节数量
var LEN_INFO_BYTE_SIZE int = 4
var (
	ST_MULTI = "multi"
	ST_SINGLE = "single"
)

var (
	ST_TCPCONN = "tcpConn"
	ST_TOKENLIB = "tokenlib"
)

var (
	CT_GET_TOKEN = "getToken"
	CT_HEARTBEAT = "heartBeat"
)

//DB_KEY
var (
	REDIS_TOKEN_KEY = "userToken"
)

