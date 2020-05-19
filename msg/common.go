// @Author: liubai
// @Date: 2020/5/7 10:08 下午
// @Desc: 公共的方法

package msg

import (
	"bufio"
	"fmt"
	"github.com/golang/protobuf/proto"
	"tcpFrame/datas/proto"
	"tcpFrame/util"
	"time"
)

//todo 根据codeType实现封装序列化sendBody的interface{}，将decoding部分脱离出去
//todo 业务自行序列化sendMsg数据，只传入一个[]byte格式的sendMsg
func SendMessage(rw *bufio.ReadWriter, cmdNo, bodyType int, sendMsg proto.Message) error {
	fmt.Println(util.RunFuncName(), cmdNo, sendMsg)

	//todo 按照codeType序列化数据
	sendHeader := &heartbeat.RequestHeader{
		CmdNo:      uint32(cmdNo),
		BodyLength: uint32(proto.Size(sendMsg)),
		BodyType:   uint32(bodyType),
		Version:    "v1.0.1",
	}
	headerBytes, _ := proto.Marshal(sendHeader)
	msgBytes, _ := proto.Marshal(sendMsg)
	bData, _ := BuildData(headerBytes, msgBytes)
	n, err := rw.Write(bData)
	err1 := rw.Flush()
	fmt.Println(util.RunFuncName(), "send data size: ", n, bData)
	time.Sleep(time.Microsecond * 10)
	if err != nil || err1 != nil {
		fmt.Println(util.RunFuncName(), "have err ", err)
		return err
	}
	return nil
}

func ReadMessage(rw *bufio.ReadWriter, headBytesChan chan []byte, msgBytesChan chan []byte, closeFlag chan int) {
	var recieveBytes []byte

	readChan := make(chan []byte, 1024)
	//从tcp iobuf中读取数据放入readChan中
	go func() {
		for {
			bData := make([]byte, 1024)
			n, err := rw.Read(bData)
			fmt.Println(util.RunFuncName(), "get data size: ", n)
			if err != nil {
				fmt.Println("链接无法读取，连接关闭。", err)
				closeFlag <- 1
				return
			}
			if n > 0 {
				bData = bData[:n]
				readChan <- bData
				fmt.Println(util.RunFuncName(), "get data: ", bData)
			}
		}
	}()

	//将上面方法读取的数据存入本地缓存recieveBytes中
	for {
		s := <-readChan
		recieveBytes = util.BytesCombine(recieveBytes, s)
		headerBytes, msgBytes, err := ParseBaseHeaderData(&recieveBytes)
		fmt.Println(util.RunFuncName(), "get rawData: ", headerBytes, msgBytes, err)

		if len(headerBytes) > 0 && len(msgBytes) > 0 {
			headBytesChan <- headerBytes
			msgBytesChan <- msgBytes
		}
	}
}
