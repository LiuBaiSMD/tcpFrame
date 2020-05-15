/*
@Author: liubai
@Date: 2020/5/14 9:07 下午
@Desc: use for what
*/

package main

import (
	"github.com/streadway/amqp"
	"time"
	"log"
	"tcpPractice/util"
	"fmt"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	util.FailOnError(err, "Failed to connect to server")
	defer conn.Close()

	ch, err := conn.Channel()
	util.FailOnError(err, "Failed to connect to channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", //name
		false,   //durable
		false,   //delete when usused
		false,   // exclusive
		false,   //no-wait
		nil,     // arguments
	)

	util.FailOnError(err, "Failed to declare a queue")
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // arguments
	)
	util.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		count:=1
		start := time.Now().Unix()
		for d := range msgs {
			count++
			log.Printf("Received a message : %s", d.Body)
			if count>1000000{
				end := time.Now().Unix()
				speed := (end-start)
				fmt.Println("consume speed: ",start , end, speed)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages, To exit press CTRL+C")
	<-forever
}
