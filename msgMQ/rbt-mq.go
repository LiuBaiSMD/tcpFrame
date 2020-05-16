/*
@Author: liubai
@Date: 2020/5/16 6:37 下午
@Desc: use for what
*/

package msgMQ

import (
	"fmt"
	"github.com/streadway/amqp"
	"sync"
	"tcpFrame/util"
)

var (
	RabbitMQMap map[string]RabbitMQ //使用serviceName， RabbitMq作为数据存储
)

// 定义RabbitMQ对象
type RabbitMQAMQP struct {
	rbtname      string
	passwd       string
	ipAddr       string
	port         int
	rbtconn      *amqp.Connection
	channel      *amqp.Channel
	producerList []Producer
	receiverList []Receiver
	exchangeMap map[string]ExchangeAMQP //交换机名称 ：ExchangeAMQP
	mu           sync.RWMutex
}
var NormalExtype = "direct"

// 定义队列交换机对象，对QuName进行声明然后绑定QuName RtKey ExName
type ExchangeAMQP struct {
	ExName string // 交换机名称
	ExType string // 交换机类型
	BindMap map[string]string //交换机绑定类型[QuName]RtKey,
}

//连接rabbitmq
func (r *RabbitMQAMQP) connect() error {
	RabbitUrl := fmt.Sprintf("amqp://%s:%s@%s:%d", r.rbtname, r.passwd, r.ipAddr, r.port)
	fmt.Println(util.RunFuncName(), "err: ", RabbitUrl)
	mqConn, err := amqp.Dial(RabbitUrl)
	if err != nil {
		fmt.Printf(util.RunFuncName(), "MQ打开链接失败:%s \n", err)
		return err
	}
	r.rbtconn = mqConn
	mqChan, err = mqConn.Channel()
	r.channel = mqChan // 赋值给RabbitMQ对象
	if err != nil {
		fmt.Printf("MQ打开管道失败:%s \n", err)
	}
	return nil
}

//绑定QuName RtKey ExName
func (r *RabbitMQAMQP) BindQueue(qName, rtKey, excName string)error{
	//首先检查连接中是否已有此交换机
	exc, ok := r.exchangeMap[excName]
	if !ok{
		exc = ExchangeAMQP{
			excName,
			NormalExtype,
			make(map[string]string),
		}
	}
	fmt.Println(util.RunFuncName(), exc)

	//找到对应的交换机
	//开始进行绑定分析
	//bm := exc.BindMap
	//if exc
	return nil
}