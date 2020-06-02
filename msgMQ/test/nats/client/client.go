/*
@Author: liubai
@Date: 2020/6/1 10:43 下午
@Desc: use for what
*/

package main

import (
	"flag"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/pborman/uuid"
	"log"
	"strconv"
	"time"
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
		subj = flag.String("subj", "test", "subject name")
	)
	flag.Parse()
	log.Println(*subj)
	startClient(*subj)

	time.Sleep(time.Second)
}

//send message to server
func startClient(subj string) {

	id := uuid.New()
	log.Println(id)
	nc.Publish(subj, []byte(id+" Sun "+strconv.Itoa(1)))
	nc.Publish(subj, []byte(id+" Rain "+strconv.Itoa(2)))
	nc.Publish(subj, []byte(id+" Fog "+strconv.Itoa(3)))
	nc.Publish(subj, []byte(id+" Cloudy "+strconv.Itoa(4)))
	fmt.Println("v:")
}

func checkErr(err error) bool {
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
