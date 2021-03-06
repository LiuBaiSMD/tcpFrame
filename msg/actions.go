/*
@Author: liubai
@Date: 2020/5/10 6:16 下午
@Desc: 处理消息的具体方法，需要使用的方法可通过registry模块进行注册
*/

package msg

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"log"
	"net"
	"strconv"
	"tcpFrame/conns"
	_const "tcpFrame/const"
	"tcpFrame/dao"
	heartbeat "tcpFrame/datas/proto"
	"tcpFrame/registry"
	"tcpFrame/util"
)

type ServerRfAddr struct {
}

func (b *ServerRfAddr) Communicate() registry.HttpWR {
	return func(conn net.Conn, BData []byte) error {
		log.Println("method:", util.RunFuncName()) //获取请求的方法
		return nil
	}
}

func (b *ServerRfAddr) HeartBeat() registry.HttpWR {
	return func(conn net.Conn, BData []byte) error {
		req := &heartbeat.HeartBeatReq{}
		err := proto.Unmarshal(BData, req)
		if err != nil {
			return err
		}
		log.Println(util.RunFuncName(), req) //获取请求的方法

		rsp := &heartbeat.HeartBeatRsp{
			UserId: req.UserId,
			Code:    200,
			Version: version,
		}
		msgByte, _ := proto.Marshal(rsp)
		SendMessage(conn, _const.ST_TCPCONN, _const.CT_HEARTBEAT, msgByte, req.UserId)
		//msgProto := &heartbeat.LoginRequest{}

		conns.FlushConnLive(int(req.UserId))
		return nil
	}
}

// 根据msgBody中的userId获取连接并发送数据
func handleNatsMsg(msg *nats.Msg) {
	hp := &heartbeat.MsgBody{}
	proto.Unmarshal(msg.Data, hp)
	conn := conns.GetConnByUId(int(hp.UserId)).GetConn()
	if conn == nil {
		log.Println(util.RunFuncName(), "nil conn!")
		return
	}
	log.Println(util.RunFuncName(), "send message ----> ", hp.UserId)
	SendMessage(conn, hp.ServerType, hp.CmdType, hp.MsgBytes, hp.UserId)
}

func checkToken(userId, token string) bool {
	rdsToken, err := dao.GetuserToken(userId)
	if rdsToken == token && err == nil {
		return true
	}
	return false
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
