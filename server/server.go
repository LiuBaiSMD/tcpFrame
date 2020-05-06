// @Author: liubai
// @Date: 2020/5/2 5:26 下午
// @Desc: use for what

package main

import (
	"fmt"
	"log"
	"net"
	"tcpPractice/msg"
)

func main() {

	addr := "127.0.0.1:8080"

	tcpAddr, err := net.ResolveTCPAddr("tcp",addr)

	if err != nil {
		log.Fatalf("net.ResovleTCPAddr fail:%s", addr)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("listen %s fail: %s", addr, err)
	} else {

		log.Println("rpc listening", addr)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("listener.Accept error:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	//根据连接的数据进行dispach

	fmt.Println("get a accept")
	//defer conn.Close()
	err := msg.ListenMessageServerBeforeLogin(conn)
	if err!=nil{
		fmt.Println("listenMessage error: ", err.Error())
	}
	fmt.Println("handlerConnection over")
}
