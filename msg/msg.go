// @Author: liubai
// @Date: 2020/5/7 10:08 下午
// @Desc: 简单实现的服务器逻辑模块

package msg

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
	configCs "tcpFrame/config/consul"
	"tcpFrame/conns"
	"tcpFrame/const"
	"tcpFrame/datas/proto"
	"tcpFrame/msgMQ/nats-mq"
	"tcpFrame/registry"
	"tcpFrame/util"
)

var version = "v1.0.1"
var register *registry.Base
var serverConfigs map[string][]configCs.ServerRegistry

var senderId string

//tcp连接后处理消息的入口，进行数据解读以及消息分发
func HandleConnection(conn net.Conn) {

	//设定一个连接超过此数量直接拒绝连接，防止导致之前的连接出错
	if conns.LenthConn() > _const.MAX_CONNS_LENGTH {
		conn.Close()
		return
	}

	//读取的数据通过chan交互
	headBytesChan := make(chan []byte, 1)
	msgBytesChan := make(chan []byte, 1)
	closeFlag := make(chan int, 1)

	//监听tcp层发送的消息
	go ReadMessage(conn, headBytesChan, msgBytesChan, closeFlag)

	//监听连接关闭信号，如果发生错误，将关闭连接
	go func() {
		<-closeFlag
		log.Println(util.RunFuncName(), "get wrong data, will close conn!")
		return
	}()

	//不断从recieveBytes读取数据解析
	for {

		//监听rawData数据
		headBytes := <-headBytesChan
		msgBytes := <-msgBytesChan

		header := &request.RequestHeader{}
		err := proto.Unmarshal(headBytes, header)
		if err != nil || header.UserId == 0 {
			//协议出错断开连接
			log.Println("get wrong header: ", string(headBytes))
			closeFlag <- 1
		}
		if header.ServerType == _const.ST_TCPCONN {

			// 单独处理token登录部分
			if header.CmdType == _const.CT_LOGIN_WITH_TOKEN {
				err := handleTokenLogin(conn, header.UserId, msgBytes)
				if err != nil {
					closeFlag <- 1
				}
			} else {
				err := dispatch(int(header.UserId), header.CmdType, msgBytes)
				if err != nil {
					log.Println(util.RunFuncName(), "data", header)
					closeFlag <- 1
				}
			}

		} else {
			//根根据header，指定的serverType部分，将数据发送到对应的nats频道

			// 加工一道，方便业务模块自行进行解析
			msgBody := ParseMsg2RbtByte(senderId, header.ServerType, header.CmdType, header.UserId, _const.MT_TCPCONN_SERVER, msgBytes)
			natsmq.Publish(header.ServerType, msgBody)
		}
	}
}

// dispatch根据cmdType进行处理
func dispatch(userId int, cmdType string, msgBytes []byte) error {
	// 首先检查是否有这个连接，没有则直接返回
	conn := conns.GetConnByUId(userId).GetConn()
	if conn == nil {
		log.Println(util.RunFuncName(), "nil conn")
		return errors.New("empty conn!")
	}
	handleFunc, ok := register.FuncRegistry[cmdType]
	if !ok {
		return errors.New("error cmdType:" + cmdType)
	}
	err := handleFunc(conn, msgBytes)
	return err
}

// 发送消息到io管道中，需要携带参数服务类型 指令类型 消息（字节格式） 发送的用户Id
func SendMessage(conn net.Conn, serverType, cmdType string, sendMsg []byte, userId int64) error {
	sendHeader := &request.RequestHeader{
		UserId:     userId,
		ServerType: serverType,
		CmdType:    cmdType,
		Version:    version,
	}
	headerBytes, _ := proto.Marshal(sendHeader)
	bData, err := BuildData(headerBytes, sendMsg)
	n, err1 := conn.Write(bData)
	//log.Println(util.RunFuncName(), userId, cmdType, n)
	if err != nil || err1 != nil {
		log.Println(util.RunFuncName(), userId, "have err ", err, err1, n)
		return err
	}
	return nil
}

//接受消息的方法，会将解读出来的消息传入两个chan 一个是此消息的消息头header 一个是消息体， 如果发生错误会有数据传入closeFlag chan进行监听
func ReadMessage(conn net.Conn, headBytesChan chan []byte, msgBytesChan chan []byte, closeFlag chan int) {
	var recieveBytes []byte

	readChan := make(chan []byte, 1024)
	//从tcp iobuf中读取数据放入readChan中
	go func() {
		for {
			bData := make([]byte, 1024)
			n, err := conn.Read(bData)
			if err != nil {
				log.Println(util.RunFuncName(), "链接无法读取，连接关闭。", err)
				closeFlag <- 1
				return
			}
			if n > 0 {
				bData = bData[:n]
				readChan <- bData
			}
		}
	}()

	//将上面方法读取的数据存入本地缓存recieveBytes中
	for {
		s := <-readChan
		recieveBytes = util.BytesCombine(recieveBytes, s)
		headerBytes, msgBytes, _ := Parse2HeaderAndMsg(&recieveBytes)
		if len(headerBytes) > 0 && len(msgBytes) > 0 {
			headBytesChan <- headerBytes
			msgBytesChan <- msgBytes
		}
	}
}
