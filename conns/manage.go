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

//定期清理超时没有发送心跳包的连接
func manageConnLive(){
	t := time.NewTicker(time.Second * 10)
	for{
		<-t.C
		nowTime := time.Now().Unix()
		for key := range cMap.connLiveMap {
			if (nowTime - cMap.connLiveMap[key]) > int64(_const.MAX_LOSE_HEARTBEAT * _const.HEARTBEAT_INTERVAL){
				fmt.Println(util.RunFuncName(), "conn dicconnect: ", key)
				//超时从队列中删除
				DelConnById(key)
			}
			fmt.Println(util.RunFuncName(), "conn: ", key)
		}
	}
}