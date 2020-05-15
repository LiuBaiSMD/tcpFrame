/*
@Author: liubai
@Date: 2020/5/14 9:07 下午
@Desc: use for what
*/

package main

import (
	"github.com/streadway/amqp"
	"fmt"
	"time"
	"tcpPractice/util"
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
	count := 1
	start := time.Now().Unix()
	for{
		count++
		err = ch.Publish(
			"",     //exchange
			 q.Name,     // routing key
			false,  //mandatory
			false, //immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body :      []byte(body),
			})
		if count >1000000{
			break
		}
		//util.FailOnError(err, "Failed to publish a message")
	}
	end := time.Now().Unix()
	speed := (end-start)
	fmt.Println("push speed: ",start, end, speed)
}