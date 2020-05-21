/*
@Author: liubai
@Date: 2020/5/14 9:07 下午
@Desc: use for what
*/

package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"tcpFrame/util"
)

var queueName1 = "hello"
var exchangeName1 = "exchangeName1"
var exchangeType1 = "direct"
var routeKey1 = "rbt.key1"

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	util.FailOnError(err, "Failed to connect to server")
	defer conn.Close()

	ch, err := conn.Channel()
	util.FailOnError(err, "Failed to connect to channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName1, //name
		false,      //durable
		false,      //delete when usused
		false,      // exclusive
		false,      //no-wait
		nil,        // arguments
	)

	util.FailOnError(err, "Failed to declare a queue")
	msgs, err := ch.Consume(
		q.Name,      // queue
		"consumer1", // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // arguments
	)
	util.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		count := 1
		for d := range msgs {
			count++
			log.Printf("Received a message : %s", d.Body)
		}
	}()

	go func() {
		q, err1 := ch.QueueDeclare(
			queueName1, //name
			false,      //durable
			false,      //delete when usused
			false,      // exclusive
			false,      //no-wait
			nil,        // arguments
		)
		fmt.Println("err1: ", err1)
		msgs, err1 := ch.Consume(
			q.Name,      // queue
			"consumer1", // consumer
			true,        // auto-ack
			false,       // exclusive
			false,       // no-local
			false,       // no-wait
			nil,         // arguments
		)
		fmt.Println("q.Name+2: ", q.Name)
		forever1 := make(chan bool)

		go func() {
			count := 1
			for d := range msgs {
				count++
				log.Printf("Received a message1 : %s", d.Body)
			}
		}()
		<-forever1
	}()

	log.Printf(" [*] Waiting for messages, To exit press CTRL+C")
	<-forever
}
