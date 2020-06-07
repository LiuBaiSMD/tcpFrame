// @Author: liubai
// @Date: 2020/5/2 5:26 下午
// @Desc: 模拟服务端，其多功能tcp服务

package main

import (
	"fmt"
	"log"
	"net"
	"tcpFrame/config/consul"
	"tcpFrame/conns"
	"tcpFrame/const"
	"tcpFrame/msg"
	"tcpFrame/server-registry"
	"tcpFrame/util"
	"time"
)

//本服务注册使用的ip和端口
var ipAddr string = "127.0.0.1"
var port int = 8080

func main() {
	go countConn()

	//先从服务器获取配置
	consulCfg, err := consul.GetConsulCfg(_const.CONSUL_IP, _const.CONSUL_PORT)
	natsCfg, err1 := consul.GetNatsCfg(_const.CONSUL_IP, _const.CONSUL_PORT)
	redisCfg, err2 := consul.GetRedisCfg(_const.CONSUL_IP, _const.CONSUL_PORT)
	if util.CheckNils(consulCfg, natsCfg, redisCfg)||util.CheckNotNils(err, err1, err2) {
		panic("config error!")
	}

	//注册服务
	server_registry.ConsulConnect(fmt.Sprintf("%s:%d", consulCfg.Ip, consulCfg.Port))
	serverId, err := server_registry.RegisterServer(ipAddr, port, _const.ST_TCPCONN, []string{"tcpConn"})
	if err != nil {
		log.Fatalln("服务注册失败： ", _const.ST_TCPCONN)
	}
	defer server_registry.DeRegistry(serverId)

	//启动服务
	msg.InitServer(serverId,
		msg.SetRedisCfg(*redisCfg),
		msg.SetConsulCfg(*consulCfg),
		msg.SetNatsCfg(*natsCfg),
	)
	addr := fmt.Sprintf("%s:%d", ipAddr, port)

	// 先检测是否已经使用或者已启动
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
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

func countConn() {
	for {
		time.Sleep(time.Second)
		log.Println(util.RunFuncName(), "conn length = ", conns.LenthConn())
	}
}


