// @Author: liubai
// @Date: 2020/5/7 10:08 下午
// @Desc: use for what

package msg

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	_const "tcpPractice/const"
	"tcpPractice/datas"
	"tcpPractice/util"
	"time"
)

func LoginForClient(conn net.Conn, cData datas.Request)(bool, error){
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	SendMessage(rw, _const.LOGIN_ACTION, cData)
	respone, err := GetMessage(rw)
	if err!=nil{
		return false, errors.New("no data")
	}
	repData, ok := respone.(datas.BaseData)
	if !ok{
		return false, errors.New("data error")
	}

	if repData.Action==_const.LOGIN_SUCCESS_ACTION{
		return true, nil
	}else {
		//return false, errors.New("login failed")
		return true, nil
	}

}

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

func ListenMessageClient(conn net.Conn, breakFlag chan int)(error){
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	for{
		respone, err := GetMessage(rw)
		if err!=nil{
			fmt.Println("ListenMessageClient: ", err)
			breakFlag<-1
			return errors.New("no data")
		}
		responeData, ok := respone.(datas.BaseData)
		if ok{
			fmt.Println("respone: ", responeData)
		}
	}
}