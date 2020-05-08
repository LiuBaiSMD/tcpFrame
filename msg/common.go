// @Author: liubai
// @Date: 2020/5/7 10:08 下午
// @Desc: use for what

package msg

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"tcpPractice/datas"
)

func SendMessage(conn net.Conn, msg interface{})error{
	bData, _ := json.Marshal(msg)
	_, err := conn.Write(bData)
	if err!=nil{
		return err
	}
	return nil
}

func GetMessage(conn net.Conn)(interface{}, error){
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
