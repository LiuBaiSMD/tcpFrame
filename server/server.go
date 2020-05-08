// @Author: liubai
// @Date: 2020/5/2 5:26 下午
// @Desc: use for what

package main

import (
	"fmt"
	"log"
	"net"
	"tcpPractice/conns"
	"tcpPractice/msg"
	"tcpPractice/test"
	"tcpPractice/util"
	"time"
)

func main() {
	go testConn()
	go testTcp.TestReconnect(conns.GetCMap())
	addr := "127.0.0.1:8080"

	tcpAddr, err := net.ResolveTCPAddr("tcp",addr)
	if err != nil {
		log.Fatalf("net.ResovleTCPAddr fail:%s", addr)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("listen %s fail: %s", addr, err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("listener.Accept error:", err)
			continue
		}

		go msg.HandleConnection(conn)
	}
}

func testConn(){
	for{
		time.Sleep(time.Second)
		conn := conns.GetConnByUId(10001)
		if conn!=nil{
			fmt.Println(util.RunFuncName(), "have conn")
			continue
		}
		fmt.Println(util.RunFuncName(), "have not conn")
	}
}


