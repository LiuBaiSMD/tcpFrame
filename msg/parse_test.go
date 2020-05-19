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
	"tcpFrame/util"
	"testing"
)

func Test_RequestMinLen(t *testing.T) {
	p := &heartbeat.LoginRequest{
		UserName:"wuxun",
		Password:"123456",
		Token:"abcdefghigjk",
		LoginType:1,
		Version:1,
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
