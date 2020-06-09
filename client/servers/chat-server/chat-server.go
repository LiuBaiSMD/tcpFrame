/*
@Author: liubai
@Date: 2020/5/31 4:54 下午
@Desc: 模拟后端服务的聊天系统
*/

package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"log"
	"strconv"
	"tcpFrame/config/consul"
	"tcpFrame/const"
	"tcpFrame/dao"
	"tcpFrame/datas/proto"
	"tcpFrame/msgMQ/nats-mq"
	sr "tcpFrame/server-registry"
	"tcpFrame/util"
)

var serverId string
var serverType = _const.ST_CHAT_ROOM

func main() {
	//初始化数据库
	sr.ConsulConnect(_const.CONSUL_URL)
	redisConfig, err := consul.GetRedisCfg(_const.CONSUL_IP, _const.CONSUL_PORT)
	natsConfig, err1 := consul.GetNatsCfg(_const.CONSUL_IP, _const.CONSUL_PORT)
	if redisConfig == nil || err != nil || natsConfig == nil || err1 != nil {
		panic("redis config err:" + err.Error())
	}
	dao.InitRedis(redisConfig.Password, fmt.Sprintf("%s:%d", redisConfig.Ip, redisConfig.Port), redisConfig.DB)
	serverName := serverType
	serverId, _ = sr.RegisterServer(
		"127.0.0.1",
		0,
		serverName,
		[]string{})
	fmt.Println("serverId", serverId, serverName)
	natsmq.Init(natsConfig.Ip, natsConfig.Port)
	// 订阅一个整个服务的频道
	go natsmq.AsyncNats(serverName, serverName, handleNatsMsg)
	// 订阅一个自己serverId专属的频道
	go natsmq.AsyncNats(serverId, serverId, handleNatsMsg)
	select {}
}

func handleNatsMsg(msg *nats.Msg) {

	msgBody := &request.MsgBody{}
	log.Println(util.RunFuncName(), "msgBody", msgBody)
	err := proto.Unmarshal(msg.Data, msgBody)
	if err != nil {
		return
	}
	revieverId := msgBody.SenderId
	// todo 根据msgBody.CmdType解析 msgBody.MsgBytes
	// todo 使用registry模块，将方法自动注册，然后通过cmdType进行自动调用
	if msgBody.CmdType == _const.CT_COMMUNICATE {
		pb := &request.CommunicateReq{}
		proto.Unmarshal(msgBody.MsgBytes, pb)
		log.Println(pb)
		revieveUId, err := strconv.ParseInt(pb.Message, 10, 64)
		log.Println("userId", string(revieveUId), revieveUId)
		rpb := &request.CommunicateRsp{
			//UserId:  revieveUId, // todo 模拟用户交流
			UserId:  10005, // todo 模拟用户交流
			Message: "[respone]:hello from " + strconv.FormatInt(pb.UserId,10),
		}
		rspBytes, _ := proto.Marshal(rpb)
		msgBody.ServerType = serverType
		msgBody.MsgType = _const.MT_NORMAL_SERVER
		msgBody.SenderId = serverId
		msgBody.MsgBytes = rspBytes
		rpbBytes, _ := proto.Marshal(msgBody)
		err = natsmq.Publish(revieverId, rpbBytes)
		if err != nil {
			fmt.Println(util.RunFuncName(), err)
		}
	}

}
