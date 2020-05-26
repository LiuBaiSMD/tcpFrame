/*
@Author: liubai
@Date: 2020/5/17 5:48 下午
@Desc: use for what
*/

package msgMQ

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"sync"
	"tcpFrame/util"
)

// 定义RabbitMQ对象
type RabbitMQAMQP struct {
	serviceId string
	rbtname   string
	passwd    string
	ipAddr    string
	port      int
	rbtconn   *amqp.Connection
	channel   *amqp.Channel
	//producerList []Producer
	//receiverList []Receiver
	msgRecieves  map[string]map[string]<-chan amqp.Delivery //监听消息的队列
	exchangeMap map[string]ExchangeAMQP                    //交换机名称 ：ExchangeAMQP
	mu          sync.RWMutex
}

//专门用来接收rabbitmq数据的连接
type RabbitMQmsg struct {
}

var mqConn *amqp.Connection
var mqChan *amqp.Channel
var NormalExtype = "direct"

// 定义队列交换机对象，对QuName进行声明然后绑定QuName RtKey ExName
type ExchangeAMQP struct {
	ExName     string              // 交换机名称
	ExType     string              // 交换机类型
	Durable    bool                //Durable
	AutoDelete bool                //AutoDelte
	Internal   bool                //Internal
	NoWait     bool                //NoWait
	BindMap    map[string][]string //交换机绑定类型[QuName]RtKeyList,
	RtKeyList  []string            //[string]RtKeyList
}

//连接rabbitmq
func (r *RabbitMQAMQP) connect() error {
	RabbitUrl := fmt.Sprintf("amqp://%s:%s@%s:%d", r.rbtname, r.passwd, r.ipAddr, r.port)
	fmt.Println(util.RunFuncName(), "RabbitUrl: ", RabbitUrl)
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

//绑定QuName RtKey ExName，绑定只能用一个交换机，但是交换机内部的交换规则可以是queueName和rtKey
func (r *RabbitMQAMQP) AddBindQueueInfo(qName, rtKey, excName string) error {
	//首先检查连接中是否已有此交换机
	exc, ok := r.exchangeMap[excName]
	if !ok {
		exc = ExchangeAMQP{
			excName,
			NormalExtype,
			true,
			false,
			false,
			true,
			make(map[string][]string),
			make([]string, 0),
		}
	}

	//添加QueueBind信息
	var flag int = 0
	for _, rtk := range exc.BindMap[qName] {
		if rtk == rtKey {
			flag = 1
		}
	}
	if flag == 0 {
		exc.BindMap[qName] = append(exc.BindMap[qName], rtKey)
	}
	flag = 0
	for _, rtk := range exc.RtKeyList {
		if rtk == rtKey {
			flag = 1
		}
	}
	if flag == 0 {
		exc.RtKeyList = append(exc.RtKeyList, rtKey)
	}
	r.exchangeMap[excName] = exc

	fmt.Println(util.RunFuncName(), exc, r.exchangeMap[excName])
	return nil
}

//进行队列绑定
func (r *RabbitMQAMQP) BindQueue(qName, rtKey, excName string) error {
	//首先判断是否已经有队列
	//q, err := r.channel.QueueDeclarePassive(qName, true, false, false, true, nil)
	_, err := r.channel.QueueDeclare(qName, true, false, false, true, nil)
	fmt.Println(util.RunFuncName(), err)
	if err != nil{
		fmt.Println(util.RunFuncName(), " err: ", err)
		// 队列不存在,声明队列
		// name:队列名称;durable:是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能;autoDelete:是否自动删除;noWait:是否非阻塞,
		// true为是,不等待RMQ返回信息;args:参数,传nil即可;exclusive:是否设置排他
		if err != nil {
			fmt.Printf("MQ注册队列失败:%s \n", err)
			return errors.New("MQ注册队列失败")
		}
	}
	fmt.Println(util.RunFuncName(), " err: ", err)
	// 用于检查交换机是否存在,已经存在不需要重复声明
	//err = r.channel.ExchangeDeclarePassive(excName, NormalExtype, true, false, false, true, nil)
	err = r.channel.ExchangeDeclare(excName, NormalExtype, true, false, false, true, nil)
	if err != nil {
		fmt.Println(util.RunFuncName(), "没有交换机，正在注册。。。")
		// 注册交换机
		// name:交换机名称,kind:交换机类型,durable:是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能;autoDelete:是否自动删除;
		// noWait:是否非阻塞, true为是,不等待RMQ返回信息;args:参数,传nil即可; internal:是否为内部
		if err != nil {
			fmt.Printf("MQ注册交换机失败:%s \n", err)
			return errors.New("MQ注册交换机失败")
		}
	}

	// 进行绑定操作
	err = r.channel.QueueBind(qName, rtKey, excName, true, nil)
	r.AddBindQueueInfo(qName, rtKey, excName)
	return nil
}

func (r *RabbitMQAMQP) Publish(excName, rtKey string, msg []byte) error {
	//首先判断是否有此交换机
	exc, ok := r.exchangeMap[excName]
	if !ok {
		return errors.New("have no exchange registry")
	}
	fmt.Println(util.RunFuncName(), "rtKey list: ", exc.RtKeyList)
	ok = false
	for _, rtk := range exc.RtKeyList {
		fmt.Println(util.RunFuncName(), excName, rtk)
		if rtk == rtKey {
			ok = true
		}
	}
	fmt.Println(util.RunFuncName(), r.exchangeMap[excName].RtKeyList)
	if !ok {
		return errors.New("have no rtKey registry")
	}

	// 发送任务消息
	err := r.channel.Publish(excName, rtKey, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        msg,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *RabbitMQAMQP) Consume() {
	ch := r.channel
	msgs, err := ch.Consume("hello", "consumer1", true, false, false, false, nil)
	if err != nil {
		fmt.Println(util.RunFuncName(), "err:", err)
	}
	forever := make(chan bool)

	go func() {
		count := 1
		for d := range msgs {
			count++
			fmt.Printf("Received a message : %s", d.Body)
			forever <- true
		}
	}()
	<-forever
}

func (r *RabbitMQAMQP) MakeConsumeMsg(qName, consumeName string) error {
	fmt.Println(util.RunFuncName(), "test")
	ch := r.channel
	_, err := r.channel.QueueDeclarePassive(qName, true, false, false, true, nil)
	if err != nil {
		fmt.Println(util.RunFuncName(), " err: ", err)
		// 队列不存在,声明队列
		// name:队列名称;durable:是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能;autoDelete:是否自动删除;noWait:是否非阻塞,
		// true为是,不等待RMQ返回信息;args:参数,传nil即可;exclusive:是否设置排他
		fmt.Println(util.RunFuncName(), "声明队列")
		_, err = r.channel.QueueDeclare(qName, true, false, false, true, nil)
		if err != nil {
			fmt.Printf("MQ注册队列失败:%s \n", err)
			return errors.New("MQ注册队列失败")
		}
	}
	msg, err := ch.Consume(qName, consumeName, true, false, false, false, nil)
	if err != nil {
		fmt.Println(util.RunFuncName(), err)
		return err
	}
	msgList, ok := r.msgRecieves[qName]
	if !ok {
		r.msgRecieves[qName] = make(map[string]<-chan amqp.Delivery)
		msgList = r.msgRecieves[qName]
	}

	_, ok = msgList[consumeName]
	if ok {
		fmt.Println("已有此consumer: " + consumeName+" , 将进行替换！")
	}
	msgList[consumeName] = msg
	r.msgRecieves[qName] = msgList
	fmt.Println(util.RunFuncName(), msgList)
	return err
}
