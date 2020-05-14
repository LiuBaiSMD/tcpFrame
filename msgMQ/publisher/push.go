/*
@Author: liubai
@Date: 2020/5/14 9:07 下午
@Desc: use for what
*/

package main

import (
	"github.com/streadway/amqp"
	"tcpPractice/util"
	"time"
)

func main(){
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	util.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	util.FailOnError(err, "Failed to open a channel")
	q, err := ch.QueueDeclare(
		"hello", //name
		false,  //durable
		false,  //delete when unused
		false,  //exclusive
		false,  //no wait
		nil,    //arguments
	)
	util.FailOnError(err, "Failed to declare q queue")

	body := "Hello1"
	for{
		time.Sleep(time.Second * 2)
		err = ch.Publish(
			"",     //exchange
			q.Name,     // routing key
			false,  //mandatory
			false, //immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body :      []byte(body),
			})

		util.FailOnError(err, "Failed to publish a message")
	}

}