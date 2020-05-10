// @Author: liubai
// @Date: 2020/5/7 10:08 下午
// @Desc: use for what

package msg

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"tcpPractice/conns"
	_const "tcpPractice/const"
	"tcpPractice/datas"
)

func HandleConnection(conn net.Conn) {
	//根据连接的数据进行dispach
	fmt.Println("get a accept")
	//defer conn.Close()
	err := ListenMessageServerBeforeLogin(conn)
	if err!=nil{
		fmt.Println("listenMessage error: ", err.Error())
	}
	fmt.Println("handlerConnection over")
}

func ListenMessageServerBeforeLogin(conn net.Conn)error{
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	respone, err := GetMessage(rw)
	if err!=nil{
		return errors.New("no data")
	}
	cData, ok := respone.(datas.BaseData)
	if ok && !CheckLogin(cData){
		fmt.Println("login failed!")
		//验证登录消息
		//返回登录失败信息
		defer conn.Close()
		respone := datas.Respone{
			Action:_const.LOGIN_FAILED_ACTION,
			Code:200,
		}
		err := SendMessage(rw, _const.LOGIN_FAILED_ACTION, respone)
		if err!=nil{
			return err
		}
	}
	//登录成功
	fmt.Println("login success!")
	respone = datas.Respone{
		Action:_const.LOGIN_SUCCESS_ACTION,
		Code:200,
	}
	err = SendMessage(rw, _const.LOGIN_SUCCESS_ACTION, respone)
	if err!=nil{
		return err
	}
	//将连接加入到conns连接池中，跳出循环，进行其他监听
	userClient := conns.NewClient(cData.UserId, conn, cData.UserId)
	conns.PushChan(cData.UserId, userClient)
	err = ListenMessageAfterLogin(cData.UserId, conn)
	return err
}

func ListenMessageAfterLogin(connId int,conn net.Conn)error{
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	//断开连接后从连接池中删除
	defer conns.DelConnById(connId)

	for{
		fmt.Println("ListenMessageAfterLogin")
		respone, err := GetMessage(rw)
		if err!=nil{
			return errors.New("no data")
		}
		err = DisPatch(conn, respone)
		if err!=nil{
			return err
		}
	}
}

//校验登录参数是否正确
func CheckLogin(cData datas.BaseData)bool{
	fmt.Println("login data", cData)
	if cData.Action != _const.LOGIN_ACTION{
		fmt.Println(cData.Action, _const.LOGIN_ACTION)
		return false
	}
	return true
}