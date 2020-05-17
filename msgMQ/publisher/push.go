/*
@Author: liubai
@Date: 2020/5/14 9:07 下午
@Desc: use for what
*/

package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"tcpFrame/util"
	"time"
)

var queueName1 = "hello"
var exchangeName1 = "exchangeName1"
var exchangeType1 = "direct"
var routeKey1 = "rbt.key1"

func main(){
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	util.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	util.FailOnError(err, "Failed to open a channel")
	_, err = ch.QueueDeclare(
		queueName1, //name
		false,  //durable
		false,  //delete when unused
		false,  //exclusive
		false,  //no wait
		nil,    //arguments
	)

	_, err = ch.QueueDeclare(
		queueName1+"2", //name
		false,  //durable
		false,  //delete when unused
		false,  //exclusive
		false,  //no wait
		nil,    //arguments
	)
	ch.QueueBind(queueName1, routeKey1+"2", exchangeName1, true, nil)
	go func(){
		count := 1
		for{
			count++
			time.Sleep(time.Second * 2)
			err = ch.Publish(
				exchangeName1,     //exchange
				routeKey1+"2",     // routing key
				false,  //mandatory
				false, //immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body :      []byte("world"),
				})
			if count >10000000{
				break
			}
			fmt.Println("send world!", exchangeName1+"2", routeKey1)
			//util.FailOnError(err, "Failed to publish a message")
		}
	}()

	util.FailOnError(err, "Failed to declare q queue")
	err =  ch.ExchangeDeclare(exchangeName1+"2", exchangeType1, true, false, false, true, nil)
	if err != nil {
		fmt.Printf("MQ注册交换机失败:%s \n", err)
		return
	}
	// 绑定任务
	err =  ch.QueueBind(queueName1+"2", routeKey1, exchangeName1+"2", true, nil)
	if err != nil {
		fmt.Printf("绑定队列失败:%s \n", err)
		return
	}
	body := "Hello"
	count := 1
	start := time.Now().Unix()
	for{
		count++
		time.Sleep(time.Second * 2)
		err = ch.Publish(
			exchangeName1+"2",     //exchange
			routeKey1,     // routing key
			false,  //mandatory
			false, //immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body :      []byte(body),
			})
		if count >10000000{
			break
		}
		//util.FailOnError(err, "Failed to publish a message")
	}
	end := time.Now().Unix()
	speed := (end-start)
	fmt.Println("push speed: ",start, end, speed)
}