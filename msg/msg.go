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
	"tcpPractice/datas"
)

func ListenMessageServer(conn net.Conn)error{
	for{
		respone, err := getMessage(conn)
		if err!=nil{
			return errors.New("no data")
		}
		cData, ok := respone.(datas.StructData)
		if ok{
			err = DisPatch(conn, cData)
			if err!=nil{
				return err
			}
		}
	}
}

func ListenMessageClient(conn net.Conn)(error){
	for{
		respone, err := getMessage(conn)
		if err!=nil{
			return errors.New("no data")
		}
		responeData, ok := respone.(datas.StructData)
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
		var cData datas.StructData
		err:=json.Unmarshal(bData[:n], &cData)
		if err != nil{
			fmt.Println("err:", err)
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
	cData, ok := data.(datas.StructData)
	if ok && cData.Action == "login"{
		fmt.Println("login", cData.Name, cData.PWD)
		respone := datas.StructData{
			Action:"loginRespone",
			Code:200,
		}
		err := SendMessage(conn, respone)
		if err!=nil{
			return err
		}
	}else{
		return errors.New(fmt.Sprintf("can not found 【%s】 action in registry!", cData.Action))
	}
	return nil
}
