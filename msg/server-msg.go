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
	"tcpPractice/util"
)

//func HandleConnection(conn net.Conn) {
//	//根据连接的数据进行dispach
//	fmt.Println("get a accept")
//	//defer conn.Close()
//	err := ListenMessageServerBeforeLogin(conn)
//	if err!=nil{
//		fmt.Println("listenMessage error: ", err.Error())
//	}
//	fmt.Println("handlerConnection over")
//}

func HandleConnection(conn net.Conn) {
	//将连接加入到conns连接池中，跳出循环，进行其他监听
	userClient := conns.NewClient(10001, conn, 10001)
	conns.PushChan(10001, userClient)
	//根据连接的数据进行dispach
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	var recieveBytes []byte
	readChan := make(chan []byte, 1024)
	done := make(chan int, 1)
	closeFlag := make(chan int, 1)
	//读取的数据通过chan交互

	//从tcp iobuf中读取数据放入readChan中
	go GetBytesFromBuf(rw, readChan, closeFlag)
	go BindBytesFromBuf(&recieveBytes, readChan, done, closeFlag)
	//不断从读取的数据中解析
	for{
		fmt.Println("get start")
		<-done
		fmt.Println("get over")
		rawData := ReadData(&recieveBytes)
		fmt.Println(util.RunFuncName(), "rawData: ", rawData)
		//rawData := 1
	}
}

//不断的从网络连接buf中获取数据
func GetBytesFromBuf(rw *bufio.ReadWriter, readChan chan []byte, closeFlag chan int){
	for{
		bData := make([]byte, 1024)
		n, err := rw.Read(bData)
		fmt.Println(util.RunFuncName(), "get data size: ", n)
		if err != nil{
			fmt.Println("链接无法读取，连接关闭。", err)
			closeFlag<-1
		}
		if n>0 {
			bData = bData[:n]
			readChan <- bData
			fmt.Println(util.RunFuncName(), "get data: ", bData)
		}
	}
}

func BindBytesFromBuf(byteStore *[]byte, readChan chan []byte, done chan int, closeFlag chan int){
	for{
		s := <- readChan
		fmt.Println(util.RunFuncName(), "get data: ", s)
		*byteStore = util.BytesCombine(*byteStore, s)
		done<-1
	}
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