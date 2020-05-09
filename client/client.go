// @Author: liubai
// @Date: 2020/5/2 5:27 下午
// @Desc: use for what

package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"tcpPractice/const"
	"tcpPractice/datas"
	"tcpPractice/msg"
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

	done = make(chan int, 1)
	connClose = make(chan int, 1)
	loginFlag, err := msg.LoginForClient(conn, loginData)
	if !loginFlag || err!=nil{
		fmt.Println("login failed: ", loginFlag, err)
		return
	}
	fmt.Println("login success: ")
	go msg.ListenMessageClient(conn, done)
	go msg.Heartbeat(userId, conn, connClose)
	<-done
	connClose <- 1
}

