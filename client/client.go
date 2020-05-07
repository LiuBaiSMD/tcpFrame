// @Author: liubai
// @Date: 2020/5/2 5:27 下午
// @Desc: use for what

package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"tcpPractice/const"
	"tcpPractice/datas"
	"tcpPractice/msg"
	"tcpPractice/util"
	"time"
)

func Open(addr string) (*bufio.ReadWriter, net.Conn, error) {
	fmt.Println("Dial " + addr)
	//conn, err := tls.Dial("tcp", addr, nil)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, nil, errors.New(err.Error() +  "Dialing "+addr+" failed")
	}
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), conn, nil
}

var userId = 10001
var done chan int
var loginData = datas.Request{
	Action:_const.LOGIN_ACTION,
	Name:"wuxun",
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
	loginFlag, err := msg.LoginForClient(conn, loginData)
	if !loginFlag || err!=nil{
		fmt.Println("login failed: ", loginFlag, err)
		return
	}
	go msg.ListenMessageClient(conn)
	go Heartbeat(userId, conn)
	<-done

}



func Heartbeat(userId int, conn net.Conn)error{
	timer := time.NewTicker(time.Second * 5)
	for{
		<- timer.C
		fmt.Println("heartbeat")
		heartbeatRequest := datas.Request{
			Action: "heartbeat",
			UserId:	userId,
		}
		bData, _ := json.Marshal(heartbeatRequest)
		_, err := conn.Write(bData)
		if err!=nil{
			fmt.Println(util.RunFuncName(), " : ", err)
			return err
		}
	}

	return nil
}