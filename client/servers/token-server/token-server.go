/*
@Author: liubai
@Date: 2020/5/31 4:54 下午
@Desc: 模拟服务集群中的token获取服务
*/

package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"strconv"
	"tcpFrame/const"
	"tcpFrame/dao"
	"tcpFrame/datas/proto"
	"tcpFrame/handle"
	"tcpFrame/msgMQ"
	natsmq "tcpFrame/msgMQ/nats-mq"
	sr "tcpFrame/server-registry"
	"tcpFrame/util"
)

var rdsConn *redis.Client

func main() {
	//初始化数据库
	rdsConn = dao.InitRedis("", "127.0.0.1:6379", 0)
	sr.ConsulConnect("localhost:8500")
	serverName := _const.ST_TOKENLIB
	sr.RegisterServer(
		"127.0.0.1",
		0,
		serverName,
		[]string{})
	defer sr.DeRegistryAll(serverName)

	//接受从rabbtmq发送过来的数据
	go natsmq.AsyncNats(serverName, "test", handleMsg)
	GetRbtMsg(serverName)
}

func GetRbtMsg(serverName string) {
	msgMQ.BindServiceQueue("server1", serverName)
	rspServerName := serverName + "res"
	msgMQ.BindServiceQueue("server1", rspServerName)
	msgMQ.AddConsumeMsg("server1", serverName, "consumer2")
	rbtMsg, err := msgMQ.GetConsumeMsgChan("server1", serverName, "consumer2")
	if err != nil || rbtMsg == nil {
		fmt.Println(util.RunFuncName(), err, "没有数据或连接!")
	} else {
		for {
			fmt.Println(util.RunFuncName(), "ready to get msg")
			message := <-rbtMsg
			fmt.Println("get message : ", message.Body)
			msgBody := &heartbeat.MsgBody{}
			err = proto.Unmarshal(message.Body, msgBody)
			if msgBody.CmdType == _const.CT_GET_TOKEN {
				pb := &heartbeat.TokenTcpRequest{}
				proto.Unmarshal(msgBody.MsgBytes, pb)
				s := strconv.FormatInt(pb.UserId, 10)
				token, err := handle.GetTokenReal(s, pb.UserName)
				fmt.Println(util.RunFuncName(), "token: ", pb, token, err)
				dao.SaveUserToken(s, token)
				rpb := &heartbeat.TokenTcpRespone{
					UserId: pb.UserId,
					Token:  token,
				}
				rpbBytes, _ := proto.Marshal(rpb)
				msgMQ.Publish2Service("server1", rspServerName, rpbBytes)
			}

		}
	}
	//获取该serverName下的所有服务节点信息
	//servicesMap, _ := server_registry.ServicesMap("serverNode")
}

func handleMsg(msg *nats.Msg) {
	msgBody := &heartbeat.MsgBody{}
	err := proto.Unmarshal(msg.Data, msgBody)
	fmt.Println(util.RunFuncName(), msgBody, err)
	pb := &heartbeat.TokenTcpRequest{}
	proto.Unmarshal(msgBody.MsgBytes, pb)
	fmt.Println(util.RunFuncName(), "token: ", pb)

}
