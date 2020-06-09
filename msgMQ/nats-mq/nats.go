/*
@Author: liubai
@Date: 2020/6/1 10:27 下午
@Desc: 使用queue模式
*/

package natsmq

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"tcpFrame/util"
	"log"
)

const (
	url = "nats://127.0.0.1:4222"
)

var (
	nc  *nats.Conn
	err error
)

func Init(ip string, port int) {
	url := fmt.Sprintf("nats://%s:%d", ip, port)
	if nc, err = nats.Connect(url); util.PanicErr(err) {

	}
}

//订阅一个服务的方法， subj为订阅的频道，workQueue为工作组，
//handle为接收到方法需要对数据进行操作的自动调用的方法
func AsyncNats(subj string, workQueue string, handle nats.MsgHandler) {
	_, err := nc.QueueSubscribe(subj, workQueue, handle)
	//_, err := nc.Subscribe(subj, handle)
	if err!=nil{
		log.Println(util.RunFuncName(), err)
	}
}

func Publish(subj string, msg []byte) error {
	if err := nc.Publish(subj, msg); !util.LogErr(err) {
		return err
	}
	return nil
}
