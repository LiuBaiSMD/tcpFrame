// @Author: liubai
// @Date: 2020/5/2 5:27 下午
// @Desc: 模拟客户端，在运行server中主函数后调用

package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	"os"
	"tcpFrame/const"
	"tcpFrame/datas/proto"
	"tcpFrame/msg"
	"tcpFrame/msgMQ"
	"tcpFrame/util"
	"time"
)

func Open(addr string) (*bufio.ReadWriter, net.Conn, error) {
	fmt.Println("Dial " + addr)
	//conn, err := tls.Dial("tcp", addr, nil)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, nil, errors.New(err.Error() + "Dialing " + addr + " failed")
	}
	return nil, conn, nil
}

var userId = int64(10001)
var done chan int
var connClose chan int

func main() {
	go testRbtAndServerRegist()
	_, conn, err := Open("127.0.0.1:8080")
	if err != nil {
		fmt.Println("dial failed:", err)
		os.Exit(1)
	}
	defer conn.Close()

	connClose = make(chan int, 1)
	go Heartbeat(userId, conn, connClose)
	<-connClose
}

func Heartbeat(userId int64, conn net.Conn, closeFlag chan int) error {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	timer := time.NewTicker(time.Second * time.Duration(_const.HEARTBEAT_INTERVAL))
	for {
		select {
		case <-timer.C:
			req := &heartbeat.HeartBeatReq{
				UserId: userId,
				Version:   "v1.0.1",
			}
			err := msg.SendMessage(rw, _const.ST_TOKENLIB, _const.CT_GET_TOKEN, req, userId)
			if err != nil {
				fmt.Println(util.RunFuncName(), " : ", err)
				closeFlag <- 1
				return err
			}
		}
	}
	return nil
}

func testRbtAndServerRegist() {
	//首先注册一个服务
	serverName := _const.ST_TOKENLIB

	msgMQ.BindServiceQueue("server1", serverName)
	msgMQ.AddConsumeMsg("server1", serverName, "consumer2")
	rbtMsg, err := msgMQ.GetConsumeMsgChan("server1", serverName, "consumer2")
	if err != nil || rbtMsg == nil {
		fmt.Println(util.RunFuncName(), err, "没有数据或连接!")
	}else{
		for{
			fmt.Println(util.RunFuncName(), "ready to get msg")
			message := <- rbtMsg
			fmt.Println("get message : ", message.Body)
			dp2 := &heartbeat.HeartBeatRsp{}
			err = proto.Unmarshal(message.Body, dp2)
			fmt.Println("dp2: ", dp2, " err: ", err)
		}
	}
	//获取该serverName下的所有服务节点信息
	//servicesMap, _ := server_registry.ServicesMap("serverNode")
}
