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

func (b *ServerRfAddr) Communicate() registry.HttpWR {
	return func(rw *bufio.ReadWriter, BData []byte) error {
		log.Log("method:", util.RunFuncName()) //获取请求的方法
		return nil
	}
}

func (b *ServerRfAddr) HeartBeat() registry.HttpWR {
	return func(rw *bufio.ReadWriter, BData []byte) error {
		req := &heartbeat.HeartBeatReq{}
		proto.Unmarshal(BData, req)
		log.Log("method:", util.RunFuncName(), req) //获取请求的方法

		rsp := &heartbeat.HeartBeatRsp{
			Code: 200,
			Version:  "login success!",
		}
		SendMessage(rw, _const.ST_TCPCONN, _const.CT_HEARTBEAT, rsp, req.UserId)
		//msgProto := &heartbeat.LoginRequest{}

		//conns.FlushConnLive(BData.UserId)
		return nil
	}
}
