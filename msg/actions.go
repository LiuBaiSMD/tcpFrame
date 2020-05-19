/*
@Author: liubai
@Date: 2020/5/10 6:16 下午
@Desc: 处理消息的具体方法，需要使用的方法可通过registry模块进行注册
*/

package msg

import (
	"bufio"
	"github.com/golang/protobuf/proto"
	"github.com/micro/go-micro/util/log"
	"tcpFrame/const"
	"tcpFrame/datas/proto"
	"tcpFrame/registry"
	"tcpFrame/util"
)

type ServerRfAddr struct {

}

func (b* ServerRfAddr)Communicate() registry.HttpWR{
	return func(rw *bufio.ReadWriter, BData proto.Message)error{
		log.Log("method:", util.RunFuncName()) //获取请求的方法
		return nil
	}
}

func (b* ServerRfAddr)HeartBeat() registry.HttpWR {
	return  func(rw *bufio.ReadWriter, BData proto.Message)error{
		log.Log("method:", util.RunFuncName()) //获取请求的方法
		rsp := &heartbeat.LoginRespone{
			Code:200,
			LoginState:1,
			Oms:"login success!",
		}
		SendMessage(rw, _const.CMD_LOGIN_REQ, _const.BT_LOGIN_REQ, rsp)
		//msgProto := &heartbeat.LoginRequest{}

		//conns.FlushConnLive(BData.UserId)
		return nil
	}
}
