// @Author: liubai
// @Date: 2020/5/7 10:08 下午
// @Desc: use for what

package msg

import (
	"bufio"
	"fmt"
	"net"
	"tcpPractice/const"
	"tcpPractice/datas"
	"tcpPractice/util"
	"time"
)

func Heartbeat(userId int, conn net.Conn, closeFlag chan int)error{
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	timer := time.NewTicker(time.Second * time.Duration(_const.HEARTBEAT_INTERVAL) )
	for{
		select {
		case <- timer.C:
			heartbeatRequest := datas.Request{
				Action: _const.HEARTBEAT_ACTION,
				//Action: "testheartbeat",
				UserId:	userId,
			}
			err := SendMessage(rw, _const.HEARTBEAT_ACTION, heartbeatRequest)
			if err!=nil{
				fmt.Println(util.RunFuncName(), " : ", err)
				return err
			}
		case <- closeFlag:
			fmt.Println(util.RunFuncName(), " ----> closeFlag")
			return nil
		}
	}
	return nil
}
