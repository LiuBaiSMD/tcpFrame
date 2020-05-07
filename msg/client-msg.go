// @Author: liubai
// @Date: 2020/5/7 10:08 下午
// @Desc: use for what

package msg

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	_const "tcpPractice/const"
	"tcpPractice/datas"
	"tcpPractice/util"
	"time"
)

func LoginForClient(conn net.Conn, cData datas.Request)(bool, error){
	bData, _ := json.Marshal(cData)
	_, err := conn.Write(bData)
	if err!=nil{
		conn.Close()
		return false, errors.New("login error")
	}
	respone, err := getMessage(conn)
	if err!=nil{
		return false, errors.New("no data")
	}
	cData, _ = respone.(datas.Request)
	if err!=nil{
		return false, err
	}

	if cData.Action==_const.LOGIN_SUCCESS_ACTION{
		return true, nil
	}else {
		return false, errors.New("login failed")
	}

}

func Heartbeat(userId int, conn net.Conn, closeFlag chan int)error{
	timer := time.NewTicker(time.Second * 5)
	for{
		select {
		case <- timer.C:
			heartbeatRequest := datas.Request{
				Action: _const.HEARTBEAT_ACTION,
				UserId:	userId,
			}
			bData, _ := json.Marshal(heartbeatRequest)
			_, err := conn.Write(bData)
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
	for{
		respone, err := getMessage(conn)
		if err!=nil{
			fmt.Println("ListenMessageClient: ", err)
			breakFlag<-1
			return errors.New("no data")
		}
		responeData, ok := respone.(datas.Request)
		if ok{
			fmt.Println("respone: ", responeData)
		}
	}
}