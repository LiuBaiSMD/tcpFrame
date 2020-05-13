/*
auth:   wuxun
date:   2020-05-08 10:24
mail:   lbwuxun@qq.com
desc:   how to use or use for what
*/

package util

import (
	"bufio"
	"fmt"
	"tcpPractice/datas"
	"tcpPractice/msg"
	"time"
	"tcpPractice/conns"
)

func TestReconnect(connMap conns.ConnMap){
	for{
		fmt.Println("---->", RunFuncName())
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
		rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
		msg.SendMessage(rw, "comunicate", transData)
		fmt.Println(RunFuncName(), "---->")
	}
}
