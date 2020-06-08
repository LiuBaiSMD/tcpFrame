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

//同时连接的最大数量
var MAX_CONNS_LENGTH = 10000

var (
	ST_MULTI  = "multi"
	ST_SINGLE = "single"
)

//已有的服务类型
var (
	ST_TCPCONN  = "tcpConn"
	ST_TOKENLIB = "tokenlib"
	ST_HTTP_TOKEN = "httpToken"
)

//服务中的cmd指令类型
var (
	CT_LOGIN_WITH_TOKEN = "TokenLogin"
	CT_GET_TOKEN = "getToken"
	CT_HEARTBEAT = "HeartBeat"
	CT_COMMUNICATE = "Communicate"
)

//服务中的msgBody类型
var (
	//tcpConn服务类型
	MT_TCPCONN_SERVER int32 = 1
	//普通服务集群中的服务类型
	MT_NORMAL_SERVER int32 = 2
)

//DB_KEY
var (
	REDIS_TOKEN_KEY = "userToken"
)


//TOKEN 认证状态
var (
	TOKEN_RIGHT int32 = 1
	TOKEN_WRONG int32 = 0
)