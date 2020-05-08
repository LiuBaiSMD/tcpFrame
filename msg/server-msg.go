// @Author: liubai
// @Date: 2020/5/7 10:08 下午
// @Desc: use for what

package msg


import (
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

func DisPatch(conn net.Conn, data interface{}) error{
	cData, ok := data.(datas.Request)
	if cData.Action==_const.HEARTBEAT_ACTION && cData.UserId>0{
		//更新连接的时间戳
		conns.FlushConnLive(cData.UserId)
	}
	SendMessage(conn, cData)
	fmt.Println("loginAfter", cData, ok)
	return nil
}


//校验登录参数是否正确
func CheckLogin(cData datas.Request)bool{
	fmt.Println("login data", cData)
	if cData.Action != _const.LOGIN_ACTION{
		fmt.Println(cData.Action, _const.LOGIN_ACTION)
		return false
	}
	if cData.Name==_const.LOGIN_AUTH && cData.PWD != ""{
		fmt.Println(cData.Name, cData.PWD)
		return true
	}
	return false
}


func ListenMessageServerBeforeLogin(conn net.Conn)error{
	respone, err := GetMessage(conn)
	if err!=nil{
		return errors.New("no data")
	}
	cData, ok := respone.(datas.Request)
	if ok && !CheckLogin(cData){
		fmt.Println("login failed!")
		//验证登录消息
		//返回登录失败信息
		defer conn.Close()
		respone := datas.Respone{
			Action:_const.LOGIN_FAILED_ACTION,
			Code:200,
		}
		err := SendMessage(conn, respone)
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
	err = SendMessage(conn, respone)
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
	//断开连接后从连接池中删除
	defer conns.DelConnById(connId)

	for{
		fmt.Println("ListenMessageAfterLogin")
		respone, err := GetMessage(conn)
		if err!=nil{
			return errors.New("no data")
		}
		cData, ok := respone.(datas.Request)
		if ok{
			err := DisPatch(conn, cData)
			if err!=nil{
				return err
			}
		}
	}
}