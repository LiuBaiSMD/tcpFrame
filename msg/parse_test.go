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

func Test_RequestMinLen(t *testing.T) {
	p := &heartbeat.LoginRequest{
		UserName:"wuxun",
		Password:"123456",
		Token:"abcdefghigjk",
		LoginType:1,
		Version:"v1.0.1",
	}
	pb, err := GetRequestByte(1, 1, "v1.1.1", p)
	fmt.Println(util.RunFuncName(), "pb: ", pb, "\nerr: ", err)
	hp := &heartbeat.RequestHeader{}
	proto.UnmarshalMerge(pb, hp)
	fmt.Println("hp: ", hp)
	pbb := pb[16:]
	bp := &heartbeat.LoginRequest{}
	proto.UnmarshalMerge(pbb, bp)
	fmt.Println("bp: ", bp, pbb)

}

func GetRequestByte(cmdNo, bodyType uint32, version string, body proto.Message) ([]byte, error){
	header := &heartbeat.RequestHeader{
		CmdNo:      cmdNo,
		BodyType:   bodyType,
		Version:    version,
		BodyLength: uint32(proto.Size(body)),
	}
	header.HeadLength = uint32(proto.Size(header))

	hb, err := proto.Marshal(header)
	db, err := proto.Marshal(body)
	fmt.Println(util.RunFuncName(),  hb)
	fmt.Println(util.RunFuncName(),  db)

	rb := util.BytesCombine(hb, db)
	fmt.Println(util.RunFuncName(), " rb: ", rb)
	return rb, err
}

var ioBuf []byte

func Test_ParseMsg(t *testing.T) {
	msgBody := &heartbeat.LoginRequest{
		UserName:"wuxun",
		Password:"123456",
		LoginType:1,
	}
	msgBytes, _ := proto.Marshal(msgBody)

	sendHeader := &heartbeat.RequestHeader{
		CmdNo:      1,
		BodyLength: uint32(proto.Size(msgBody)),
		BodyType:   1,
		Version:    "v1.0.1",
	}
	headerBytes, _ := proto.Marshal(sendHeader)

	ioBuf, _ = msg.BuildData(headerBytes, msgBytes)
	//加入两个结构体，模拟粘包
	ioBuf1, _ := msg.BuildData(headerBytes, msgBytes)
	for i:=0;i<len(ioBuf1);i++{
		ioBuf = append(ioBuf, ioBuf1[i])
	}

	msg.IoBuf = ioBuf
	fmt.Println(util.RunFuncName(), ioBuf)
	codeType, bRawData, err :=msg.Parse2HeaderData(&ioBuf)
	fmt.Println(codeType, bRawData, err)
	//故意加入一个字符串进行解析，
	ioBuf = util.BytesCombine(ioBuf, []byte("hello world!"))
	codeType, bRawData, err =msg.Parse2HeaderData(&ioBuf)
	fmt.Println(codeType, bRawData, err)
	fmt.Println(string(ioBuf))
	codeType, bRawData, err =msg.Parse2HeaderData(&ioBuf)
	fmt.Println(codeType, bRawData, err)

}