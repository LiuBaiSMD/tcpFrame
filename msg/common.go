// @Author: liubai
// @Date: 2020/5/7 10:08 下午
// @Desc: use for what

package msg

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"tcpPractice/datas"
	"tcpPractice/util"
	"time"
)

func SendMessage(rw *bufio.ReadWriter, action string, sendMsg interface{})error{
	fmt.Println(util.RunFuncName(), action, sendMsg)
	RBdata, _ := json.Marshal(sendMsg)
	sendBody := datas.BaseData{
		Action: action,
		BData: RBdata,
		UserId:10001,
	}
	//bData, _ := json.Marshal(sendBody)
	//再加工一次，
	bData, _ := BuildData(1, sendBody)
	n, err := rw.Write(bData)
	rw.Write(bData)
	fmt.Println(util.RunFuncName(), "send data size: ", n)
	err1 := rw.Flush()
	time.Sleep(time.Microsecond*10)
	fmt.Println(util.RunFuncName(), "rw flush")
	if err!=nil||err1!=nil{
		fmt.Println(util.RunFuncName(), "have err ", err)
		return err
	}
	return nil
}

func GetMessage(rw *bufio.ReadWriter)(interface{}, error){
	fmt.Println(util.RunFuncName(), "start")
	bData := make([]byte, 1024)
	n, err := rw.Read(bData)
	fmt.Println(util.RunFuncName(), "get data size: ", n)
	if err != nil{
		fmt.Println("链接无法读取，连接关闭。", err)
		return nil, errors.New("链接无法读取，连接关闭。")
	}
	if n>0 {
		var cData datas.BaseData
		err:=json.Unmarshal(bData[:n], &cData)
		if err!=nil{
			fmt.Println(util.RunFuncName(), "recieve data: ", string(bData))
			return nil, err
		}
		fmt.Println(util.RunFuncName(), cData.Action, cData)
		return cData, nil
	}
	return nil, errors.New("no data")
}
