/*
auth:   wuxun
date:   2020-05-08 20:23
mail:   lbwuxun@qq.com
desc:   how to use or use for what
*/

package conns

import "fmt"

//定期清理超时没有发送心跳包的连接
func manageConnLive(){
	for key := range cMap.connLiveMap {
		fmt.Println(key, " : ", cMap.connLiveMap[key])
	}
}