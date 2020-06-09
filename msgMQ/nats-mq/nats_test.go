/*
@Author: liubai
@Date: 2020/6/1 10:27 下午
@Desc: use for what
*/

package natsmq_test

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"tcpFrame/msgMQ/nats-mq"
	"testing"
	"time"
)


const (
	url = "nats://localhost:4222"
	subj = "weather"
)

var (
	nc  *nats.Conn
	err error
)
func init() {
	if nc, err = nats.Connect(url); checkErr(err) {

	}
}
func Test_rabbitMq(t *testing.T) {

	startServer(subj, "s1")
	startServer(subj, "s2")
	startServer(subj, "s3")
	//wait for subscribe complete
	time.Sleep(1 * time.Second)

	startClient(subj)

	select {}
}


func TestAsyncNats(t *testing.T) {
	natsmq.Init("127.0.0.1", 4222)
	natsmq.AsyncNats("test", "test", handleMsg)
	natsmq.AsyncNats("test", "test", handleMsg)
	//natsmq.AsyncNats(*servername, *queueGroup+"test", handleMsg)

	err := natsmq.Publish("test", []byte("hello world!"))
	fmt.Println("test:", err)
	time.Sleep(time.Second)
}
func handleMsg(msg *nats.Msg) {
	log.Println("Received a message From Async : ", string(msg.Data))
}


//send message to server
func startClient(subj string) {
	nc.Publish(subj, []byte("Sun"))
	nc.Publish(subj, []byte("Rain"))
	nc.Publish(subj, []byte("Fog"))
	nc.Publish(subj, []byte("Cloudy"))
}

//receive message
func startServer(subj, name string) {
	go sync(nc, subj, name)
	go async(nc, subj, name)
}

func async(nc *nats.Conn, subj, name string) {
	nc.Subscribe(subj, func(msg *nats.Msg) {
		fmt.Println(name, "Received a message From Async : ", string(msg.Data))
	})
}

func sync(nc *nats.Conn, subj, name string) {
	subscription, err := nc.SubscribeSync(subj)
	checkErr(err)
	if msg, err := subscription.NextMsg(10 * time.Second); checkErr(err) {
		if msg != nil {
			fmt.Println(name, "Received a message From Sync : ", string(msg.Data))
		}
	} else {
		//handle timeout
	}

}

func checkErr(err error) bool {
	if err != nil {
		log.Println("nats 连接失败！")
		log.Fatal(err)
		return false
	}
	return true
}

