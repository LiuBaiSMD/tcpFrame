// @Author: liubai
// @Date: 2020/5/2 5:26 下午
// @Desc: 模拟服务端，其多功能tcp服务

package main

import (
	"bufio"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
	"tcpFrame/conns"
	"tcpFrame/const"
	"tcpFrame/datas/proto"
	"tcpFrame/msg"
	"tcpFrame/server-registry"
	"tcpFrame/util"
	"time"
)

func main() {
	go testConn()
	server_registry.ConsulConnect("localhost:8500")
	serverId, err := server_registry.RegisterServer("127.0.0.1", 8080, _const.ST_TCPCONN, []string{"tcpConn"})
	if err != nil{
		log.Fatalln("服务注册失败： ", _const.ST_TCPCONN)
	}
	defer server_registry.DeRegistry(serverId)
	msg.InitServer(serverId)
	//go testTcp.TestReconnect(conns.GetCMap())
	addr := "127.0.0.1:8080"

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatalf("net.ResovleTCPAddr fail:%s", addr)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("listen %s fail: %s", addr, err)
	}
	fmt.Println(util.RunFuncName(), err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("listener.Accept error:", err)
			continue
		}
		go msg.HandleConnection(conn)
	}
}

func testConn() {
	for {
		time.Sleep(time.Second)
		conn := conns.GetConnByUId(10001)
		if conn != nil {
			fmt.Println(util.RunFuncName(), "have conn")
			continue
		}
		fmt.Println(util.RunFuncName(), "have not conn , conn lengt= ", conns.LenthConn())
	}
}

func TestReconnect(connMap conns.ConnMap) {
	for {
		fmt.Println("---->", util.RunFuncName())
		time.Sleep(time.Second * 3)
		connClinet := conns.GetConnByUId(10001)
		if connClinet == nil {
			continue
		}
		conn := connClinet.GetConn()
		rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
		rsp := &heartbeat.LoginRespone{
			UserId:     0,
			Code:       200,
			LoginState: 1,
			Oms:        "login success!",
		}
		msgByte, _ := proto.Marshal(rsp)
		msg.SendMessage(rw, _const.ST_TCPCONN, _const.CT_HEARTBEAT, msgByte, 0)
		fmt.Println(util.RunFuncName(), "---->")
	}
}
