// @Author: liubai
// @Date: 2020/5/7 10:08 下午
// @Desc: 简单实现的服务器逻辑模块

package msg

import (
	"bufio"
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	"tcpFrame/conns"
	"tcpFrame/const"
	"tcpFrame/datas/proto"
	"tcpFrame/registry"
	"tcpFrame/util"
	"time"
)
var register *registry.Base
func init(){
	var rfaddr1 ServerRfAddr
	register = registry.Registery(&rfaddr1)
}

//tcp连接后处理消息的入口，进行数据解读以及消息分发
func HandleConnection(conn net.Conn) {

	//在登录成功后，将conn加入到conns连接池中,进行其他行为监听，
	//先模拟用户userId为100001的连接进入
	// todo 此部分将移入用户登录模块中
	userClient := conns.NewClient(10001, conn, 10001)
	conns.PushChan(10001, userClient)

	//读取的数据通过chan交互
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	headBytesChan := make(chan []byte, 1)
	msgBytesChan := make(chan []byte, 1)
	closeFlag := make(chan int, 1)


	//监听tcp层发送的消息
	go ReadMessage(rw, headBytesChan, msgBytesChan, closeFlag)

	//监听连接关闭信号，如果发生错误，将关闭连接
	go func(){
		<-closeFlag
		fmt.Println(util.RunFuncName(), "get wrong data, will close conn!")
		conn.Close()
		return
	}()

	//不断从recieveBytes读取数据解析
	for{

		//监听rawData数据
		headBytes := <- headBytesChan
		msgBytes := <- msgBytesChan
		fmt.Println(util.RunFuncName(), "will encode Data ", headBytes, msgBytes)

		header := &heartbeat.RequestHeader{}
		err := proto.Unmarshal(headBytes, header)
		if err!=nil{
			//协议出错断开连接
			fmt.Println("get wrong header: ", string(headBytes))
			closeFlag<-1
		}
		fmt.Println(util.RunFuncName(), "get header: ", header)
		//todo 根据codeType实现反序列化bRawData的interface{}，将encoding部分脱离出去
		if int(header.CmdNo)==_const.CMD_LOGIN_REQ{
			msg := &heartbeat.LoginRequest{}
			err := proto.Unmarshal(msgBytes, msg)
			if err !=nil{
				//协议出错断开连接
				fmt.Println("get wrong rawData: ", string(msgBytes))
				closeFlag<-1
			}
			fmt.Println(util.RunFuncName(), "proto: ", msg)
		}else if int(header.CmdNo)==_const.CMD_HEARTBEAT{
			msg := &heartbeat.LoginRequest{}
			err := proto.Unmarshal(msgBytes, msg)
			if err !=nil{
				//协议出错断开连接
				fmt.Println("get wrong rawData: ", string(msgBytes))
				closeFlag<-1
			}
			fmt.Println(util.RunFuncName(), "proto: ", msg)
		}
		//todo 通过codeType解析数据，进行dispatch
	}
}

//todo 根据codeType实现封装序列化sendBody的interface{}，将decoding部分脱离出去
//todo 业务自行序列化sendMsg数据，只传入一个[]byte格式的sendMsg
func SendMessage(rw *bufio.ReadWriter, cmdNo, bodyType int, sendMsg proto.Message) error {
	fmt.Println(util.RunFuncName(), cmdNo, sendMsg)

	//todo 按照codeType序列化数据
	sendHeader := &heartbeat.RequestHeader{
		CmdNo:      uint32(cmdNo),
		BodyLength: uint32(proto.Size(sendMsg)),
		BodyType:   uint32(bodyType),
		Version:    "v1.0.1",
	}
	headerBytes, _ := proto.Marshal(sendHeader)
	msgBytes, _ := proto.Marshal(sendMsg)
	bData, _ := BuildData(headerBytes, msgBytes)
	n, err := rw.Write(bData)
	err1 := rw.Flush()
	fmt.Println(util.RunFuncName(), "send data size: ", n, bData)
	time.Sleep(time.Microsecond * 10)
	if err != nil || err1 != nil {
		fmt.Println(util.RunFuncName(), "have err ", err)
		return err
	}
	return nil
}

func ReadMessage(rw *bufio.ReadWriter, headBytesChan chan []byte, msgBytesChan chan []byte, closeFlag chan int) {
	var recieveBytes []byte

	readChan := make(chan []byte, 1024)
	//从tcp iobuf中读取数据放入readChan中
	go func() {
		for {
			bData := make([]byte, 1024)
			n, err := rw.Read(bData)
			fmt.Println(util.RunFuncName(), "get data size: ", n)
			if err != nil {
				fmt.Println("链接无法读取，连接关闭。", err)
				closeFlag <- 1
				return
			}
			if n > 0 {
				bData = bData[:n]
				readChan <- bData
				fmt.Println(util.RunFuncName(), "get data: ", bData)
			}
		}
	}()

	//将上面方法读取的数据存入本地缓存recieveBytes中
	for {
		s := <-readChan
		recieveBytes = util.BytesCombine(recieveBytes, s)
		headerBytes, msgBytes, err := Parse2HeaderAndMsg(&recieveBytes)
		fmt.Println(util.RunFuncName(), "get rawData: ", headerBytes, msgBytes, err)

		if len(headerBytes) > 0 && len(msgBytes) > 0 {
			headBytesChan <- headerBytes
			msgBytesChan <- msgBytes
		}
	}
}