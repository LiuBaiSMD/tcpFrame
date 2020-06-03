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
	"tcpFrame/msgMQ/nats-mq"
	sr "tcpFrame/server-registry"
	"tcpFrame/util"
)

var rdsConn *redis.Client

var serverId string

func main() {
	//初始化数据库
	rdsConn = dao.InitRedis("", "127.0.0.1:6379", 0)
	sr.ConsulConnect("localhost:8500")
	serverName := _const.ST_TOKENLIB
	serverId, _ = sr.RegisterServer(
		"127.0.0.1",
		0,
		serverName,
		[]string{})
	fmt.Println("serverId", serverId)
	//接受从rabbtmq发送过来的数据
	go natsmq.AsyncNats(serverName, "test", handleMsg)
	select {}
}

func handleMsg(msg *nats.Msg) {
	msgBody := &heartbeat.MsgBody{}
	err := proto.Unmarshal(msg.Data, msgBody)
	revieverId := msgBody.SenderId
	// todo 根据msgBody.CmdType解析 msgBody.MsgBytes
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
	rspBytes, _ := proto.Marshal(rpb)
	msgBody.SenderId = serverId
	msgBody.MsgBytes = rspBytes
	rpbBytes, _ := proto.Marshal(msgBody)
	err = natsmq.Publish(revieverId, rpbBytes)
	if err!=nil{
		fmt.Println(util.RunFuncName(), err)
	}
}
