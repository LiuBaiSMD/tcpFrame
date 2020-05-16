/*
@Author: liubai
@Date: 2020/5/14 9:39 下午
@Desc: use for what
*/

package msgMQ_test

import (
	"fmt"
	"tcpFrame/msgMQ"
	"testing"
)

type TestPro struct {
	msgContent   string
}

func Test_RabbitMq(t *testing.T){
	msg := fmt.Sprintf("这是测试任务")
	tp := &TestPro{
		msg,
	}
	queueExchange := &msgMQ.QueueExchange{
		"test.rabbit",
		"rabbit.key",
		"test.rabbit.mq",
		"direct",
	}
	mq := msgMQ.New(queueExchange)
	mq.RegisterProducer(tp)
	mq.RegisterReceiver(tp)
	mq.Start()
	select {

	}
}


// 实现发送者
func (t *TestPro) MsgContent() string {
	return t.msgContent
}

// 实现接收者
func (t *TestPro) Consumer(dataByte []byte) error {
	fmt.Println(string(dataByte))
	return nil
}

