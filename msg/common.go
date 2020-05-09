// @Author: liubai
// @Date: 2020/5/7 10:08 下午
// @Desc: use for what

package msg

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"tcpPractice/datas"
	"tcpPractice/util"
	"bufio"
)

func SendMessage(conn net.Conn, action string, sendMsg interface{})error{
	fmt.Println(util.RunFuncName(), action, sendMsg)
	RBdata, _ := json.Marshal(sendMsg)
	sendBody := datas.BaseData{
		Action: action,
		BData: RBdata,
		UserId:10001,
	}
	bData, _ := json.Marshal(sendBody)
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	_, err := rw.Write(bData)
	rw.Flush()
	if err!=nil{
		fmt.Println(util.RunFuncName(), "have err ", err)
		return err
	}
	return nil
}

func GetMessage(conn net.Conn)(interface{}, error){
	fmt.Println(util.RunFuncName(), "start")
	bData := make([]byte, 512)
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	n, err := rw.Read(bData)
	fmt.Println(util.RunFuncName(), "over")
	if err != nil{
		fmt.Println("链接无法读取，连接关闭。", err)
		return nil, errors.New("链接无法读取，连接关闭。")
	}
	if n>0 {
		var cData datas.BaseData
		err:=json.Unmarshal(bData[:n], &cData)
		if err!=nil{
			return nil, err
		}
		fmt.Println(util.RunFuncName(), cData.Action, cData)
		return cData, nil
	}
	return nil, errors.New("no data")
}
