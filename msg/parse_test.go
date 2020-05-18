/*
@Author: liubai
@Date: 2020/5/18 11:32 下午
@Desc: use for what
*/

package msg_test

import (
	"fmt"
	protoM "github.com/gogo/protobuf/proto"
	"tcpFrame/datas/proto"
	"tcpFrame/msg"
	"tcpFrame/util"
	"testing"
)

func Test_RequestMinLen(t *testing.T) {
	fmt.Println(util.RunFuncName(), msg.RequestMinLen())
	req := &heartbeat.RequestHeader{}
	req.CmdNo = 1
	req.BodyLength = 1
	req.BodyType = 1
	req.HeadLength = 1
	req.Version = "1.0.1"
	b, _ := protoM.Marshal(req)
	fmt.Println(util.RunFuncName(), b, string(b[9:]))
}
