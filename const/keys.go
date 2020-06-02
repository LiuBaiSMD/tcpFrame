/*
auth:   wuxun
date:   2020-06-02 13:01
mail:   lbwuxun@qq.com
desc:   get key in const rule
*/

package _const

var (
	REQ_REAR = "_req"
	RSP_REAR = "_rsp"
)

func GetServerReqKey(serverName string) string {
	return serverName+REQ_REAR
}

func GetServerRspKey(serverName string) string {
	return serverName+RSP_REAR
}