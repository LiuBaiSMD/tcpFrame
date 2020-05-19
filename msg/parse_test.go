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
	req := &heartbeat.RequestHeader{
		CmdNo: 2,
	}
	req.BodyLength = 1
	req.BodyType = 1
	req.HeadLength = 1
	req.Version = "1.0.112"
	b, _ := proto.Marshal(req)
	fmt.Println(util.RunFuncName(), b, string(b[9:]))
	var pb = &heartbeat.RequestHeader{}
	twob := util.BytesCombine(b, b)
	//直接解析即可，取出解析后的size
	err := proto.Unmarshal(twob, pb)
	fmt.Println("Unmarshal: ", pb, " \nerr: ", err, twob, proto.Size(pb))

	twob = twob[len(twob)/2:]
	err = proto.UnmarshalMerge(twob, pb)
	fmt.Println("Unmarshal: ", pb, " \nerr: ", err, twob, pb)
	GetRequestByte(1, 1, "v1.1.1", pb)


}
func GetRequestByte(cmdNo, bodyType uint32, version string, body proto.Message) {
	header := &heartbeat.RequestHeader{
		CmdNo:      cmdNo,
		BodyType:   bodyType,
		Version:    version,
		BodyLength: uint32(proto.Size(body)),
	}
	header.HeadLength = uint32(proto.Size(header))

	hb, err := proto.Marshal(header)
	db, err := proto.Marshal(body)
	fmt.Println(util.RunFuncName(), db, err, header, hb)
	rb := util.BytesCombine(hb, db)
	fmt.Println(util.RunFuncName(), " rb: ", rb)
}
