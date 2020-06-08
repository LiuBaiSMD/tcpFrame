/*
auth:   wuxun
date:   2020-05-08 20:23
mail:   lbwuxun@qq.com
desc:   how to use or use for what
*/

package conns

import (
	"fmt"
	_const "tcpFrame/const"
	"tcpFrame/util"
	"time"
)

var nowTime int64
//定期清理超时没有发送心跳包的连接
func manageConnLive(){
	t := time.NewTicker(time.Second * 10)
	for{
		<-t.C
		nowTime = time.Now().Unix()
		cMap.connLiveMap.Range(delTimeOutConn)
	}
}

func delTimeOutConn(k, v interface{}) bool{
	if (nowTime - v.(int64)) > int64(_const.MAX_LOSE_HEARTBEAT * _const.HEARTBEAT_INTERVAL){
		fmt.Println(util.RunFuncName(), "conn dicconnect: ", k)
		//超时从队列中删除
		DelConnById(k.(int))
	}
	fmt.Println(util.RunFuncName(), "conn: ", k)
	return true
}