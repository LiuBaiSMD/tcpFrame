// @Author: liubai
// @Date: 2020/5/7 10:08 下午
// @Desc: 简单实现的服务器逻辑模块

package msg

import (
	"bufio"
	"fmt"
	"encoding/json"
	"net"
	"tcpFrame/conns"
	"tcpFrame/datas"
	"tcpFrame/util"
)


//tcp连接后处理消息的入口，进行数据解读以及消息分发
func HandleConnection(conn net.Conn) {

	//在登录成功后，将conn加入到conns连接池中,进行其他行为监听，
	//先模拟用户userId为100001的连接进入
	// todo 此部分将移入用户登录模块中
	userClient := conns.NewClient(10001, conn, 10001)
	conns.PushChan(10001, userClient)

	//读取的数据通过chan交互
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	codeTypeChan := make(chan int, 1)
	bRawChan := make(chan []byte, 1)
	closeFlag := make(chan int, 1)


	//监听tcp层发送的消息
	go ReadMessage(rw, codeTypeChan, bRawChan, closeFlag)

	//监听连接关闭信号，如果发生错误，将关闭连接
	go func(){
		<-closeFlag
		fmt.Println(util.RunFuncName(), "get wrong data, will close conn!")
		conn.Close()
		return
	}()

	//不断从recieveBytes读取数据解析
	for{

		//监听rawData数据
		codeType := <- codeTypeChan
		bRawData := <- bRawChan
		fmt.Println(util.RunFuncName(), "will encode Data ", codeType, bRawData)

		//todo 根据codeType实现反序列化bRawData的interface{}，将encoding部分脱离出去
		if codeType==1 && len(bRawData)>0{
			var rawData datas.BaseData
			err := json.Unmarshal(bRawData, &rawData)
			if err !=nil{
				//协议出错断开连接
				fmt.Println("get wrong rawData: ", string(bRawData))
				closeFlag<-1
			}
		}
		//todo 通过codeType解析数据，进行dispatch
	}
}
