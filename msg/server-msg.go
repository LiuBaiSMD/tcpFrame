// @Author: liubai
// @Date: 2020/5/7 10:08 下午
// @Desc: use for what

package msg

import (
	"bufio"
	"fmt"
	"net"
	"tcpPractice/conns"
	"tcpPractice/util"
)

func HandleConnection(conn net.Conn) {

	//在登录成功后，将conn加入到conns连接池中,进行其他行为监听，
	//先模拟用户userId为100001的连接进入
	// todo 此部分将移入用户登录模块中
	userClient := conns.NewClient(10001, conn, 10001)
	conns.PushChan(10001, userClient)


	//根据连接的数据进行dispach
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	var recieveBytes []byte
	readChan := make(chan []byte, 1024)
	//done := make(chan int, 1)
	closeFlag := make(chan int, 1)
	//读取的数据通过chan交互

	//从tcp iobuf中读取数据放入readChan中
	go func(){
		for{
			bData := make([]byte, 1024)
			n, err := rw.Read(bData)
			fmt.Println(util.RunFuncName(), "get data size: ", n)
			if err != nil{
				fmt.Println("链接无法读取，连接关闭。", err)
				closeFlag<-1
				return
			}
			if n>0 {
				bData = bData[:n]
				readChan <- bData
				fmt.Println(util.RunFuncName(), "get data: ", bData)
			}
		}
	}()

	//将上面方法读取的数据存入本地缓存recieveBytes中
	// todo 改进部分，不需要通过done管道驱动，
	go func(){
		for{
			s := <- readChan
			recieveBytes = util.BytesCombine(recieveBytes, s)
			codeType, bRawData, err := ReadData(&recieveBytes)
			if codeType!=0 && len(bRawData) > 0{
				fmt.Println(util.RunFuncName(), "rawData: ", codeType, bRawData, err )
			}
		}
	}()

	//监听连接关闭信号，准备关闭连接
	go func(){
		<-closeFlag
		conn.Close()
		return
	}()
	//go BindBytesFromBuf(&recieveBytes, readChan, done, closeFlag)
	//不断从recieveBytes读取数据解析
	select{

	}
	//for{
	//	<-done
	//	err := ReadData(&recieveBytes)
	//	fmt.Println(util.RunFuncName(), "rawData: ", recieveBytes, err)
	//}
}

