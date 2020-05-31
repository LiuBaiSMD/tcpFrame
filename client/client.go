// @Author: liubai
// @Date: 2020/5/2 5:27 下午
// @Desc: 模拟客户端，在运行server中主函数后调用

package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"tcpFrame/const"
	"tcpFrame/datas/proto"
	"tcpFrame/msg"
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
var userName = "wuxun"
var done chan int
var connClose chan int

func main() {
	//go testRbtAndServerRegist()
	_, conn, err := Open("127.0.0.1:8080")
	if err != nil {
		fmt.Println("dial failed:", err)
		os.Exit(1)
	}
	defer conn.Close()

	connClose = make(chan int, 1)
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	GetToken(rw, userId, userName)
	//go Heartbeat(userId, rw, connClose)
	<-connClose
}

func Heartbeat(userId int64, rw *bufio.ReadWriter, closeFlag chan int) error {
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

func GetToken(rw *bufio.ReadWriter, userId int64, userName string) error {
	timer := time.NewTicker(time.Second * time.Duration(_const.HEARTBEAT_INTERVAL))
	for {
		select {
		case <-timer.C:
			req := &heartbeat.TokenTcpRequest{
				UserId: userId,
				UserName: userName,
				Version:   "v1.0.1",
			}
			msg.SendMessage(rw, _const.ST_TOKENLIB, _const.CT_GET_TOKEN, req, userId)
			// 获取一个token


		}
	}
	return nil
}