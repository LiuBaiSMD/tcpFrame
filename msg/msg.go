/*
auth:   wuxun
date:   2020-05-04 14:47
mail:   lbwuxun@qq.com
desc:   how to use or use for what
*/

package msg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"errors"
	"net"
	"tcpPractice/conns"
	"tcpPractice/datas"
)

func ListenMessageServerBeforeLogin(conn net.Conn)error{
	var loginCount = 0
	for{
		loginCount+=1
		respone, err := getMessage(conn)
		if err!=nil{
			return errors.New("no data")
		}
		cData, ok := respone.(datas.Request)
		if ok{
			//验证登录消息
			if CheckLogin(cData){
				//登录成功
				//将连接加入到conns连接池中，跳出循环，进行其他监听
				userClient := conns.NewClient(cData.UserId, conn, cData.UserId)
				conns.PushChan(cData.UserId, userClient)
				break
			}
			//返回登录失败信息
			respone := datas.Respone{
				Action:"loginFailed",
				Code:200,
			}
			err := SendMessage(conn, respone)
			if err!=nil{
				return err
			}
		}
		if loginCount>=3{
			conn.Close()
			return nil
		}
	}
	ListenMessageAfterLogin(conn)
	return nil
}

func ListenMessageAfterLogin(conn net.Conn)error{
	for{
		respone, err := getMessage(conn)
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

//校验登录参数是否正确
func CheckLogin(cData datas.Request)bool{
	if cData.Action != "login"{
		return false
	}
	if cData.Name=="wuxun" && cData.PWD != ""{
		return true
	}
	return false
}

func ListenMessageClient(conn net.Conn)(error){
	for{
		respone, err := getMessage(conn)
		if err!=nil{
			return errors.New("no data")
		}
		responeData, ok := respone.(datas.Request)
		if ok{
			fmt.Println("respone: ", responeData)
		}
	}
}

func getMessage(conn net.Conn)(interface{}, error){
	bData := make([]byte, 128)
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	n, err := rw.Read(bData)
	if err != nil{
		fmt.Println("链接无法读取，连接关闭。", err)
		return nil, errors.New("链接无法读取，连接关闭。")
	}
	if n>0 {
		var cData datas.Request
		err:=json.Unmarshal(bData[:n], &cData)
		if err != nil{
			fmt.Println("err:", err, string(bData))
		}
		if err!=nil{
			return nil, err
		}
		return cData, nil
	}
	return nil, errors.New("no data")
}

func SendMessage(conn net.Conn, msg interface{})error{
	bData, _ := json.Marshal(msg)
	_, err := conn.Write(bData)
	if err!=nil{
		return err
	}
	return nil
}

func DisPatch(conn net.Conn, data interface{}) error{
	cData, ok := data.(datas.Request)
	fmt.Println("loginAfter", cData, ok)
	return nil
}
