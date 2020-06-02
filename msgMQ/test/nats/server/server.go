/*
@Author: liubai
@Date: 2020/6/1 10:43 下午
@Desc: use for what
*/

package main

import (
	"flag"
	"github.com/nats-io/nats.go"
	"log"
)

const (
	//url   = "nats://192.168.3.125:4222"
	url = nats.DefaultURL
)

var (
	nc  *nats.Conn
	err error
)

func init() {
	if nc, err = nats.Connect(url); checkErr(err) {
		//
	}
}

func main() {
	var (
		servername = flag.String("servername", "y", "name for server")
		queueGroup = flag.String("group", "q", "group name for Subscribe")
		subj       = flag.String("subj", "test", "subject name")
	)
	flag.Parse()

	log.Println(*servername, *queueGroup, *subj)
	startService(*subj, *servername+" worker1", *queueGroup)
	startService(*subj, *servername+" worker2", *queueGroup)
	startService(*subj, *servername+" worker3", *queueGroup)

	select {}
}

//receive message
func startService(subj, name, queue string) {
	go async(nc, subj, name, queue)
}

func async(nc *nats.Conn, subj, name, queue string) {
	nc.QueueSubscribe(subj, queue, handleMsg)
}

func handleMsg(msg *nats.Msg) {
	log.Println("Received a message From Async : ", string(msg.Data))
}

func checkErr(err error) bool {
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

