/*
@Author: liubai
@Date: 2020/5/31 4:54 下午
@Desc: 模拟服务集群中的token获取服务
*/

package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	_const "tcpFrame/const"
	heartbeat "tcpFrame/datas/proto"
	"tcpFrame/handle"
	"tcpFrame/msgMQ"
	sr "tcpFrame/server-registry"
	"tcpFrame/util"
)

func main() {
	serverName := _const.ST_TOKENLIB
	sr.RegisterServer(
		"127.0.0.1",
		0,
		serverName,
		[]string{})
	defer sr.DeRegistryAll(serverName)

	//接受从rabbtmq发送过来的数据
	GetRbtMsg(serverName)
}

func GetRbtMsg(serverName string) {
	msgMQ.BindServiceQueue("server1", serverName)
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
			if msgBody.CmdType == _const.CT_GET_TOKEN{
				pb := &heartbeat.TokenTcpRequest{}
				proto.Unmarshal(msgBody.MsgBytes, pb)
				token, err := handle.GetTokenReal(string(pb.UserId), pb.UserName)
				fmt.Println(util.RunFuncName(), "token: ", token, err)
			}

		}
	}
	//获取该serverName下的所有服务节点信息
	//servicesMap, _ := server_registry.ServicesMap("serverNode")
}
