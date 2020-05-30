/*
@Author: liubai
@Date: 2020/5/17 5:47 下午
@Desc: use for what
*/

package msgMQ_test

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"tcpFrame/datas/proto"
	"tcpFrame/msgMQ"
	"tcpFrame/util"
	"testing"
)

func Test_rabbitMq(t *testing.T) {
	err := msgMQ.BindServiceQueue("server1", "exchangeName1")
	fmt.Println(util.RunFuncName(), err)
	dp := &heartbeat.LoginRequest{
		UserName:"wuxun",
		Password:"123456",
		LoginType:1,
	}
	db, _ := proto.Marshal(dp)
	err = msgMQ.Publish2Service("server1", "exchangeName1",  db)
	fmt.Println(util.RunFuncName(), err)
	err = msgMQ.AddConsumeMsg("server1", "hello", "consumer1")
	fmt.Println(util.RunFuncName(), err)
	rbtMsg, err := msgMQ.GetConsumeMsgChan("server1", "hello", "consumer1")
	if err != nil || rbtMsg == nil {
		fmt.Println(util.RunFuncName(), err, "没有数据或连接!")
	}else{
		message := <- rbtMsg
		fmt.Println("get message : ", message.Body)
		dp2 := &heartbeat.LoginRequest{}
		err = proto.Unmarshal(message.Body, dp2)
		fmt.Println("dp2: ", dp2, " err: ", err)
	}
}
