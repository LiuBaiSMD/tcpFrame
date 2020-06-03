/*
@Author: liubai
@Date: 2020/5/18 11:32 下午
@Desc: use for what
*/

package msg_test

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"tcpFrame/datas/proto"
	"tcpFrame/msg"
	"tcpFrame/util"
	"testing"
)

var ioBuf []byte

func Test_ParseMsg(t *testing.T) {
	msgBody := &heartbeat.LoginRequest{
		UserName:  "wuxun",
		Password:  "123456",
		LoginType: 1,
	}
	msgBytes, _ := proto.Marshal(msgBody)

	sendHeader := &heartbeat.RequestHeader{
		CmdType: "getToken",
		Version: "v1.0.1",
	}
	headerBytes, _ := proto.Marshal(sendHeader)

	ioBuf, _ = msg.BuildData(headerBytes, msgBytes)
	//加入两个结构体，模拟粘包
	ioBuf1, _ := msg.BuildData(headerBytes, msgBytes)
	for i := 0; i < len(ioBuf1); i++ {
		ioBuf = append(ioBuf, ioBuf1[i])
	}

	msg.IoBuf = ioBuf
	fmt.Println(util.RunFuncName(), ioBuf)
	codeType, bRawData, err := msg.Parse2HeaderAndMsg(&ioBuf)
	fmt.Println(codeType, bRawData, err)
	//故意加入一个字符串进行解析，
	ioBuf = util.BytesCombine(ioBuf, []byte("hello world!"))
	codeType, bRawData, err = msg.Parse2HeaderAndMsg(&ioBuf)
	fmt.Println(codeType, bRawData, err)
	fmt.Println(string(ioBuf))
	codeType, bRawData, err = msg.Parse2HeaderAndMsg(&ioBuf)
	fmt.Println(codeType, bRawData, err)

}

func Test_protoChange(t *testing.T) {
	protoMsg := &heartbeat.LoginRequest{
		UserId:   1,
		UserName: "wuxun",
		Password: "123456",
		Version:  "v1.1.1",
	}
	changeProto(protoMsg)

}

func changeProto(msgProto proto.Message) {
	pb, err := proto.Marshal(msgProto)
	if err != nil {
		fmt.Println("err: ", err)
	}
	hp := &heartbeat.LoginRequest{}
	err = proto.Unmarshal(pb, hp)
	fmt.Println(util.RunFuncName(), "err: ", err, "msgProto: ", hp, "\nbinary: ", pb)
	fmt.Println(util.BytesToBinaryString(pb))

}

func Test_ParseMsg2RbtByte(t *testing.T) {
	dp := &heartbeat.TokenTcpRequest{
		UserId:   10001,
		UserName: "wuxun",
	}
	pb, _ := proto.Marshal(dp)
	db := msg.ParseMsg2RbtByte("test", "token", pb)
	container := &heartbeat.TokenTcpRequest{}
	parseBytes2Pb(db, container)
	fmt.Println(util.RunFuncName(), container)
}

func parseBytes2Pb(db []byte, container proto.Message) error {
	msgBody := &heartbeat.MsgBody{}
	err := proto.Unmarshal(db, msgBody)
	if err != nil {
		return err
	}
	err = proto.Unmarshal(msgBody.MsgBytes, container)
	if err != nil {
		return err
	}
	fmt.Println(util.RunFuncName(), container)
	return nil
}
