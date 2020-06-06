// @Author: liubai
// @Date: 2020/5/7 10:08 下午
// @Desc: 简单实现的服务器逻辑模块

package msg

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"log"
	"net"
	configCs "tcpFrame/config/consul"
	"tcpFrame/conns"
	"tcpFrame/const"
	"tcpFrame/dao"
	"tcpFrame/datas/proto"
	"tcpFrame/msgMQ/nats-mq"
	"tcpFrame/registry"
	"tcpFrame/util"
	"strconv"
)

var version = "v1.0.1"
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
		log.Println(err)
	}
	serverConfigs = make(map[string][]configCs.ServerRegistry)
	mJ := make([]configCs.ServerRegistry, 1)
	json.Unmarshal(multiConfig, &mJ)
	serverConfigs[_const.ST_MULTI] = mJ

	//消息中间件订阅
	natsmq.AsyncNats(serverId, serverId, handleNatsMsg)
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
		log.Println(util.RunFuncName(), "get wrong data, will close conn!")
		conn.Close()
		return
	}()

	//不断从recieveBytes读取数据解析
	for {

		//监听rawData数据
		headBytes := <-headBytesChan
		msgBytes := <-msgBytesChan

		header := &heartbeat.RequestHeader{}
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
					closeFlag <- 1
				}
			}

		} else {
			//根根据header，指定的serverType部分，将数据发送到对应的nats频道
			serverName := header.ServerType

			// 加工一道，方便业务模块自行进行解析
			msgBody := ParseMsg2RbtByte(senderId, header.CmdType, header.UserId, _const.MT_TCPCONN_SERVER, msgBytes)
			natsmq.Publish(serverName, msgBody)
		}
	}
}

// 发送消息到io管道中，需要携带参数服务类型 指令类型 消息（字节格式） 发送的用户Id
func SendMessage(rw *bufio.ReadWriter, serverType, cmdType string, sendMsg []byte, userId int64) error {
	log.Println(util.RunFuncName(), serverType, sendMsg)

	sendHeader := &heartbeat.RequestHeader{
		UserId:     userId,
		ServerType: serverType,
		CmdType:    cmdType,
		Version:    version,
	}
	headerBytes, _ := proto.Marshal(sendHeader)
	bData, _ := BuildData(headerBytes, sendMsg)
	_, err := rw.Write(bData)
	err1 := rw.Flush()
	if err != nil || err1 != nil {
		log.Println(util.RunFuncName(), "have err ", err)
		return err
	}
	return nil
}

//接受消息的方法，会将解读出来的消息传入两个chan 一个是此消息的消息头header 一个是消息体， 如果发生错误会有数据传入closeFlag chan进行监听
func ReadMessage(rw *bufio.ReadWriter, headBytesChan chan []byte, msgBytesChan chan []byte, closeFlag chan int) {
	var recieveBytes []byte

	readChan := make(chan []byte, 1024)
	//从tcp iobuf中读取数据放入readChan中
	go func() {
		for {
			bData := make([]byte, 1024)
			n, err := rw.Read(bData)
			if err != nil {
				log.Println("链接无法读取，连接关闭。", err)
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

// 根据msgBody中的userId获取连接并发送数据
func handleNatsMsg(msg *nats.Msg) {
	hp := &heartbeat.MsgBody{}
	proto.Unmarshal(msg.Data, hp)
	rw := conns.GetConnByUId(int(hp.UserId)).GetRwBuf()
	if rw == nil {
		log.Println(util.RunFuncName(), "nil conn!")
		return
	}
	SendMessage(rw, _const.ST_TOKENLIB, hp.CmdType, hp.MsgBytes, hp.UserId)
}

func checkToken(userId, token string) bool {
	rdsToken, err := dao.GetuserToken(userId)
	if rdsToken == token && err == nil {
		return true
	}
	return false
}

// dispatch根据cmdType进行处理
func dispatch(userId int, cmdType string, msgBytes []byte) error {
	// 首先检查是否有这个连接，没有则直接返回
	rw := conns.GetConnByUId(userId).GetRwBuf()
	if rw == nil {
		log.Println(util.RunFuncName(), "nil conn")
		return errors.New("empty conn!")
	}
	handleFunc, ok := register.FuncRegistry[cmdType]
	if !ok {
		return errors.New("error cmdType:" + cmdType)
	}
	err := handleFunc(rw, msgBytes)
	return err
}

func handleTokenLogin(conn net.Conn, userId int64, msgBytes []byte) error {
	msg := &heartbeat.TokenTcpRequest{}
	err := proto.Unmarshal(msgBytes, msg)
	if err != nil {
		//协议出错断开连接
		log.Println("get wrong rawData: ", string(msgBytes))
		return errors.New("get wrong rawData: " + string(msgBytes))
	}
	if checkToken(strconv.FormatInt(userId, 10), msg.Password) {
		// 登录成功
		log.Println("认证成功！", userId)
		userClient := conns.NewClient(int(userId), conn, int(msg.UserId))
		conns.PushChan(int(userId), userClient)
	} else {
		log.Println("认证失败！", userId)
		return errors.New("认证失败:" + string(userId))
	}
	return nil
}
