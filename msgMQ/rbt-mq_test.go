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
	err:=msgMQ.RabbitMq.BindQueue("hello", "rbt.key1", "exchangeName1")
	fmt.Println(util.RunFuncName(), err)
	err=msgMQ.RabbitMq.Publish("exchangeName1", "rbt.key1", []byte("hello world!"))
	fmt.Println(util.RunFuncName(), err)
}
