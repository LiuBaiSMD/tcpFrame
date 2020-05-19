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
	"tcpFrame/datas"
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
		return nil, nil, errors.New(err.Error() +  "Dialing "+addr+" failed")
	}
	return nil, conn, nil
}

var userId = 10001
var done chan int
var connClose chan int
var loginData = datas.Request{
	Action:_const.LOGIN_ACTION,
	Name:_const.LOGIN_AUTH,
	PWD:"123456",
	UserId:userId,
}

func main() {
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

func Heartbeat(userId int, conn net.Conn, closeFlag chan int)error{
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	timer := time.NewTicker(time.Second * time.Duration(_const.HEARTBEAT_INTERVAL) )
	for{
		select {
		case <- timer.C:
			req := &heartbeat.LoginRequest{
				UserName:"wuxun",
				Password:"123456",
				Token:"abcdefghigjk",
				LoginType:1,
				Version:1,
			}
			err := msg.SendMessage(rw, _const.CMD_HEARTBEAT, _const.BT_LOGIN_REQ, req)
			if err!=nil{
				fmt.Println(util.RunFuncName(), " : ", err)
				closeFlag<-1
				return err
			}
		}
	}
	return nil
}
