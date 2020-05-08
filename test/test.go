/*
auth:   wuxun
date:   2020-05-08 10:24
mail:   lbwuxun@qq.com
desc:   how to use or use for what
*/

package testTcp

import (
	"fmt"
	"tcpPractice/datas"
	"tcpPractice/util"
	"tcpPractice/msg"
	"time"
	"tcpPractice/conns"
)

func TestReconnect(connMap conns.ConnMap){
	for{
		fmt.Println("---->", util.RunFuncName())
		time.Sleep(time.Second * 3)
		connClinet := conns.GetConnByUId(10001)
		if connClinet == nil{
			continue
		}
		conn := connClinet.GetConn()
		transData := datas.Request{
			Action:"comunicate",
			Name:"testReconnect",
		}
		msg.SendMessage(conn, transData)
		fmt.Println(util.RunFuncName(), "---->")
	}
}
