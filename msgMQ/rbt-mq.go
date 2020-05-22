/*
@Author: liubai
@Date: 2020/5/16 6:37 下午
@Desc: use for what
*/

package msgMQ

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"tcpFrame/util"
)

var (
	RabbitMQMap map[string]*RabbitMQAMQP //使用serviceName， RabbitMq作为数据存储
)

func init() {
	RabbitMQMap = make(map[string]*RabbitMQAMQP)
	err := NewRabbitMQAMQP("server1", "guest", "guest", "localhost", 5672)
	if err != nil {
		panic(err)
	}
}

func NewRabbitMQAMQP(rbtServiceName, rbtname, passwd, ipAddr string, port int) error {
	rabbitAMQP := RabbitMQAMQP{
		serviceId:   rbtServiceName,
		rbtname:     rbtname,
		passwd:      passwd,
		ipAddr:      ipAddr,
		port:        port,
		exchangeMap: make(map[string]ExchangeAMQP),
		msgRecieves: make(map[string]map[string]<-chan amqp.Delivery),
	}
	err := rabbitAMQP.connect()
	RabbitMQMap[rabbitAMQP.serviceId] = &rabbitAMQP
	return err
}

func Publish2Service(serviceId, excName, routeKey string, msgBytes []byte) error {
	// 先判断是否有这个服务的消息机
	rbtmq, ok := RabbitMQMap[serviceId]
	if !ok {
		return errors.New("have no serviceId rabbitMq: " + serviceId)
	}

	err := rbtmq.Publish(excName, routeKey, msgBytes)
	return err
}

func BindServiceQueue(serviceId, excName, qName, rtKey string) error {
	rbtmq, ok := RabbitMQMap[serviceId]
	if !ok {
		return errors.New("have no serviceId rabbitMq: " + serviceId)
	}
	fmt.Println(util.RunFuncName(), qName, rtKey, excName)
	err := rbtmq.BindQueue(qName, rtKey, excName)
	return err
}

func AddConsumeMsg(serviceId, qName, consumeName string) error {
	rbtmq, ok := RabbitMQMap[serviceId]
	if !ok {
		return errors.New("have no serviceId rabbitMq: " + serviceId)
	}
	err := rbtmq.MakeConsumeMsg(qName, consumeName)
	return err
}

func GetConsumeMsgChan(serviceId, qName, consumeName string) (<-chan amqp.Delivery, error) {
	rbtmq, ok := RabbitMQMap[serviceId]
	if !ok {
		return nil, errors.New("have no serviceId rabbitMq: " + serviceId)
	}
	msgChanList, ok := rbtmq.msgRecieves[qName]
	if ok {
		msgChan, ok := msgChanList[consumeName]
		fmt.Println("ok: ", ok)
		return (msgChan), nil
	}
	return nil, errors.New("get no msgChan!")
}
