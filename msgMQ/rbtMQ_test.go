/*
@Author: liubai
@Date: 2020/5/17 5:47 下午
@Desc: use for what
*/

package msgMQ_test

import (
	"fmt"
	"tcpFrame/msgMQ"
	"tcpFrame/util"
	"testing"
)

func Test_rabbitMq(t *testing.T) {
	err := msgMQ.BindServiceQueue("server1", "exchangeName1", "hello", "rbt.key1")
	fmt.Println(util.RunFuncName(), err)
	err = msgMQ.Publish2Service("server1", "exchangeName1", "rbt.key1", []byte("hello world!"))
	fmt.Println(util.RunFuncName(), err)
	msgMQ.RabbitMQMap["server1"].Consume()
	fmt.Println(util.RunFuncName(), err)
}
