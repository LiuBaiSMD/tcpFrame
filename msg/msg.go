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
	bData := make([]byte, 128)
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	for {
		n, err := rw.Read(bData)
		if err != nil{
			fmt.Println("链接无法读取，连接关闭。", err)
			return nil
		}
		if n>0{
			var cData datas.StructData
			err:=json.Unmarshal(bData[:n], &cData)
			if err != nil{
				fmt.Println("err:", err)
			}
			fmt.Println("respone: ", cData, cData.N, cData.S)
			err = DisPatch(conn, cData)
			if err!=nil{
				return err
			}
		}
	}
	return nil
}

func ListenMessageClient(conn net.Conn)error{
	bData := make([]byte, 128)
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	for {
		n, err := rw.Read(bData)
		if err != nil{
			fmt.Println("链接无法读取，连接关闭。", err)
			return nil
		}
		if n>0{
			var cData datas.StructData
			err:=json.Unmarshal(bData[:n], &cData)
			if err != nil{
				fmt.Println("err:", err)
			}
			fmt.Println("respone: ", string(bData[:n]), cData)
			if err!=nil{
				return err
			}
		}
	}
	return nil
}

func SendMessage(conn net.Conn, msg interface{})error{
	bData, _ := json.Marshal(msg)
	_, err := conn.Write(bData)
	if err!=nil{
		return err
	}
	return nil
}

func DisPatch(conn net.Conn, data datas.StructData) error{
	if data.Action == "login"{
		fmt.Println("login", data.Name, data.PWD)
		respone := datas.StructData{
			Action:"loginRespone",
			Code:200,
		}
		err := SendMessage(conn, respone)
		if err!=nil{
			return err
		}
	}else{
		return errors.New(fmt.Sprintf("can not found 【%s】 action in registry!", data.Action))
	}
	return nil
}
