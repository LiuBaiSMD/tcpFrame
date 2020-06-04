// @Author: liubai
// @Date: 2020/5/7 10:08 下午
// @Desc: 简单实现的服务器逻辑模块

package msg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"net"
	"strconv"
	configCs "tcpFrame/config/consul"
	"tcpFrame/conns"
	"tcpFrame/const"
	"tcpFrame/dao"
	"tcpFrame/datas/proto"
	"tcpFrame/msgMQ/nats-mq"
	"tcpFrame/registry"
	"tcpFrame/util"
)

var version = "1.0.1"
var register *registry.Base
var serverConfigs map[string][]configCs.ServerRegistry

var senderId string

//tcp连接服注册方法
func InitServer(serverId string) {
	//初始化数据库
	dao.InitRedis("", "127.0.0.1:6379", 0)
	senderId = serverId
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

	//消息中间件订阅
	natsmq.AsyncNats(serverId, serverId, testHandle)
}

//tcp连接后处理消息的入口，进行数据解读以及消息分发
func HandleConnection(conn net.Conn) {
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
		//fmt.Println(util.RunFuncName(), "will encode Data ", headBytes, msgBytes)

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
			fmt.Println("---->tcpConn", header)
			if header.CmdType == _const.CT_LOGIN_WITH_TOKEN{
				msg := &heartbeat.TokenTcpRequest{}
				err := proto.Unmarshal(msgBytes, msg)
				if err != nil {
					//协议出错断开连接
					fmt.Println("get wrong rawData: ", string(msgBytes))
					closeFlag <- 1
				}
				fmt.Println(util.RunFuncName(), "proto: ", msg)
				if checkToken(strconv.FormatInt(msg.UserId, 10), msg.Password){
					// 登录成功
					// todo 此部分将移入用户登录模块中
					fmt.Println("认证成功！", msg.UserId)
					userClient := conns.NewClient(int(msg.UserId), conn, int(msg.UserId))
					conns.PushChan(int(msg.UserId), userClient)
				}else{
					fmt.Println("认证失败！", msg.UserId)
					closeFlag<-1
				}
			}

		} else {
			serverName := header.ServerType

			// 加工一道，方便业务模块自行进行解析
			msgBody := ParseMsg2RbtByte(senderId, header.CmdType, header.UserId, _const.MT_TCPCONN_SERVER, msgBytes)
			natsmq.Publish(serverName, msgBody)
		}
	}
}

func SendMessage(rw *bufio.ReadWriter, serverType, cmdType string, sendMsg []byte, userId int64) error {
	fmt.Println(util.RunFuncName(), serverType, sendMsg)

	//todo 按照codeType序列化数据
	sendHeader := &heartbeat.RequestHeader{
		UserId:     userId,
		ServerType: serverType,
		CmdType:    cmdType,
		Version:    "v1.0.1",
	}
	headerBytes, _ := proto.Marshal(sendHeader)
	bData, _ := BuildData(headerBytes, sendMsg)
	_, err := rw.Write(bData)
	err1 := rw.Flush()
	//fmt.Println(util.RunFuncName(), "send data size: ", n, bData)
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
			//fmt.Println(util.RunFuncName(), "get data size: ", n)
			if err != nil {
				fmt.Println("链接无法读取，连接关闭。", err)
				closeFlag <- 1
				return
			}
			if n > 0 {
				bData = bData[:n]
				readChan <- bData
				//fmt.Println(util.RunFuncName(), "get data: ", bData)
			}
		}
	}()

	//将上面方法读取的数据存入本地缓存recieveBytes中
	for {
		s := <-readChan
		recieveBytes = util.BytesCombine(recieveBytes, s)
		headerBytes, msgBytes, _ := Parse2HeaderAndMsg(&recieveBytes)
		//fmt.Println(util.RunFuncName(), "get rawData: ", headerBytes, msgBytes, err)

		if len(headerBytes) > 0 && len(msgBytes) > 0 {
			headBytesChan <- headerBytes
			msgBytesChan <- msgBytes
		}
	}
}

func testHandle(msg *nats.Msg){
	hp := &heartbeat.MsgBody{}
	proto.Unmarshal(msg.Data, hp)
	//fmt.Println(util.RunFuncName(), "hp:", hp)

	// todo 根据cmdType解析数据，以及在msgBody中添加serverType
	pb := &heartbeat.TokenTcpRespone{}
	if pb.Result == _const.TOKEN_RIGHT{

	}
	proto.Unmarshal(hp.MsgBytes, pb)
	rw := conns.GetConnByUId(int(pb.UserId)).GetRwBuf()
	if rw==nil{
		fmt.Println(util.RunFuncName(), "nil conn!")
		return
	}
	msgByte, _ := proto.Marshal(pb)
	SendMessage(rw, _const.ST_TOKENLIB, hp.CmdType, msgByte, hp.UserId)
}

func checkToken(userId, token string) bool {
	rdsToken, err := dao.GetuserToken(userId)
	if rdsToken == token && err == nil {
		return true
	}
	return false
}
