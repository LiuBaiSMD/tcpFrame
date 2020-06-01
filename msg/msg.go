// @Author: liubai
// @Date: 2020/5/7 10:08 下午
// @Desc: 简单实现的服务器逻辑模块

package msg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	configCs "tcpFrame/config/consul"
	"tcpFrame/conns"
	"tcpFrame/const"
	"tcpFrame/datas/proto"
	"tcpFrame/msgMQ"
	"tcpFrame/registry"
	"tcpFrame/util"
	"time"
)

var register *registry.Base
var serverConfigs map[string][]configCs.ServerRegistry

func init() {
	var rfaddr1 ServerRfAddr
	register = registry.Registery(&rfaddr1)

	// todo 读取config配置
	multiConfig, err := configCs.ReaderConfig("127.0.0.1", 8500, []string{"serverRegistry", _const.ST_MULTI})
	if err != nil {
		fmt.Println(err)
	}
	serverConfigs = make(map[string][]configCs.ServerRegistry)
	mJ := make([]configCs.ServerRegistry, 1)
	json.Unmarshal(multiConfig, &mJ)
	serverConfigs[_const.ST_MULTI] = mJ

	for _, cfg := range (serverConfigs[_const.ST_MULTI]) {
		msgMQ.BindServiceQueue("server1", cfg.Name)
	}

}

//tcp连接后处理消息的入口，进行数据解读以及消息分发
func HandleConnection(conn net.Conn) {
	//test 模块
	go testRspToken()
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
	go func() {
		<-closeFlag
		fmt.Println(util.RunFuncName(), "get wrong data, will close conn!")
		conn.Close()
		return
	}()

	//不断从recieveBytes读取数据解析
	for {

		//监听rawData数据
		headBytes := <-headBytesChan
		msgBytes := <-msgBytesChan
		fmt.Println(util.RunFuncName(), "will encode Data ", headBytes, msgBytes)

		header := &heartbeat.RequestHeader{}
		err := proto.Unmarshal(headBytes, header)
		if err != nil {
			//协议出错断开连接
			fmt.Println("get wrong header: ", string(headBytes))
			closeFlag <- 1
		}
		fmt.Println(util.RunFuncName(), "get header: ", header)
		//todo 根根据header部分，将数据发送到对应的rabbitmq
		if header.ServerType == _const.ST_TCPCONN {
			msg := &heartbeat.LoginRequest{}
			err := proto.Unmarshal(msgBytes, msg)
			if err != nil {
				//协议出错断开连接
				fmt.Println("get wrong rawData: ", string(msgBytes))
				closeFlag <- 1
			}
			fmt.Println(util.RunFuncName(), "proto: ", msg)
		} else {
			serverName := header.ServerType
			msgBody := ParstMsg2RbtByte(header.CmdType, msgBytes)
			msgMQ.Publish2Service("server1", serverName, msgBody)
		}
	}
}

//todo 根据codeType实现封装序列化sendBody的interface{}，将decoding部分脱离出去
//todo 业务自行序列化sendMsg数据，只传入一个[]byte格式的sendMsg
func SendMessage(rw *bufio.ReadWriter, serverType, cmdType string, sendMsg proto.Message, userId int64) error {
	fmt.Println(util.RunFuncName(), serverType, sendMsg)

	//todo 按照codeType序列化数据
	sendHeader := &heartbeat.RequestHeader{
		UserId:     userId,
		ServerType: serverType,
		BodyLength: uint32(proto.Size(sendMsg)),
		CmdType:    cmdType,
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

func testRspToken() {
	serverName := _const.ST_TOKENLIB
	rspServerName := serverName + "res"
	msgMQ.BindServiceQueue("server1", rspServerName)
	msgMQ.AddConsumeMsg("server1", rspServerName, "consumer2")
	rbtMsg, err := msgMQ.GetConsumeMsgChan("server1", rspServerName, "consumer2")
	if err != nil || rbtMsg == nil {
		fmt.Println(util.RunFuncName(), err, "没有数据或连接!")
	} else {
		for {
			message := <-rbtMsg
			pb := &heartbeat.TokenTcpRespone{}
			proto.Unmarshal(message.Body, pb)
			conn := conns.GetConnByUId(int(pb.UserId)).GetConn()
			if conn==nil{
				fmt.Println(util.RunFuncName(), "nil conn!")
				continue
			}
			rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
			SendMessage(rw, _const.ST_TOKENLIB, _const.CT_GET_TOKEN, pb, 10001)
			fmt.Println(util.RunFuncName(), "send: ", pb)
		}

	}
}

