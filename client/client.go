// @Author: liubai
// @Date: 2020/5/2 5:27 下午
// @Desc: 模拟客户端，在运行server中主函数后调用

package main

import (
	"bufio"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
	"net/http"
	"strconv"
	"tcpFrame/const"
	"tcpFrame/datas/proto"
	"tcpFrame/msg"
	"tcpFrame/util"
	"time"
)

var startUserId = int64(10001)
var userName = "wuxun"
var done chan int

func main() {
	for u := startUserId; u < 11001; u++ {
		go testClient(u)
		fmt.Println(util.RunFuncName(), u)
		time.Sleep(time.Microsecond * 10)
	}
	select {}
}

func testClient(userId int64) {
	var connClose chan int
	//go testRbtAndServerRegist()

	//首先通过http请求获取token
	token := httpGetToken(strconv.FormatInt(userId, 10), userName)
	if token == "" {
		log.Fatal("token 获取失败！")
	}
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("dial failed:", err)
		//os.Exit(1)
		return
	}
	defer conn.Close()

	go func() {
		rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
		headBytesChan := make(chan []byte, 1)
		msgBytesChan := make(chan []byte, 1)
		closeFlag := make(chan int, 1)

		//监听tcp层发送的消息
		go msg.ReadMessage(rw, headBytesChan, msgBytesChan, closeFlag)
		for {
			headerBytes := <-headBytesChan
			msgBytes := <-msgBytesChan
			hp := &request.RequestHeader{}
			mp := &request.TokenTcpRespone{}
			proto.Unmarshal(headerBytes, hp)
			proto.Unmarshal(msgBytes, mp)
			fmt.Println(util.RunFuncName(), hp)
			fmt.Println(util.RunFuncName(), mp)
		}
	}()

	connClose = make(chan int, 100)
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	//var st chan string
	loginWithToken(rw, userId, userName, token)
	go Heartbeat(userId, rw, connClose)
	go Chat(userId, rw, connClose)
	<-connClose
	return
	//os.Exit(2)
}

func Heartbeat(userId int64, rw *bufio.ReadWriter, closeFlag chan int) error {
	timer := time.NewTicker(time.Second * time.Duration(_const.HEARTBEAT_INTERVAL))
	for {
		select {
		case <-timer.C:
			req := &request.HeartBeatReq{
				UserId:  userId,
				Version: "v1.0.1",
			}
			msgByte, _ := proto.Marshal(req)
			err := msg.SendMessage(rw, _const.ST_TCPCONN, _const.CT_HEARTBEAT, msgByte, userId)
			if err != nil {
				fmt.Println(util.RunFuncName(), " : ", err)
				closeFlag <- 1
				return err
			}
		}
	}
	return nil
}

func Chat(userId int64, rw *bufio.ReadWriter, closeFlag chan int) error {
	timer := time.NewTicker(time.Second * time.Duration(_const.HEARTBEAT_INTERVAL))
	for {
		select {
		case <-timer.C:
			time.Sleep(time.Second)
			req := &request.CommunicateReq{
				UserId:  userId,
				//UserId:  10005,
				Message: strconv.FormatInt(userId-1, 10), // todo 模拟给上一个用户发消息
				Version: "v1.0.1",
			}
			msgByte, _ := proto.Marshal(req)
			err := msg.SendMessage(rw, _const.ST_CHAT_ROOM, _const.CT_COMMUNICATE, msgByte, userId)
			if err != nil {
				fmt.Println(util.RunFuncName(), " : ", err)
				//closeFlag <- 1
				//return err
			}
		}
	}
	return nil
}

func loginWithToken(rw *bufio.ReadWriter, userId int64, userName, token string) error {
	req := &request.TokenTcpRequest{
		UserId:   userId,
		UserName: userName,
		Password: token,
		Version:  "v1.1.1",
	}
	msgByte, _ := proto.Marshal(req)
	msg.SendMessage(rw, _const.ST_TCPCONN, _const.CT_LOGIN_WITH_TOKEN, msgByte, userId)
	return nil
}

func httpGetToken(userId, userName string) string {
	// 请求token
	client := &http.Client{}

	//生成要访问的url
	url := fmt.Sprintf("http://127.0.0.1:8081/getToken?userId=%s&userName=%s", userId, userName)
	//提交请求
	reqest, err := http.NewRequest("GET", url, nil)

	if err != nil {
		panic(err)
	}

	//处理返回结果
	response, _ := client.Do(reqest)
	tokenBData := make([]byte, 1024)
	n, _ := response.Body.Read(tokenBData)
	fmt.Println(userId, string(tokenBData[:n]), err)
	if n > 0 && err == nil {
		return string(tokenBData[:n])
	}
	return ""
}
