// @Author: liubai
// @Date: 2020/5/2 5:27 下午
// @Desc: use for what

package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

type ComplexData struct{
	N int `json:"N"`
	S  string `json:"S"`
}

func Open(addr string) (*bufio.ReadWriter, net.Conn, error) {
	// Dial the remote process.
	// Note that the local port is chosen on the fly. If the local port
	// must be a specific one, use DialTCP() instead.
	fmt.Println("Dial " + addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, nil, errors.New(err.Error() +  "Dialing "+addr+" failed")
	}
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), conn, nil
}

func main() {
	rw, conn, err := Open("127.0.0.1:8080")
	//conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("dial failed:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("rw", rw)
	var cData = ComplexData{
		N:15,
		S:"wuxun",
	}
	//readBuffer := make([]byte, 512)
	//writeBuffer := []byte("i am connector!")
	for{
		bData, _ := json.Marshal(cData)
		time.Sleep(time.Second * 5)
		n, err := rw.Write(bData)
		if err != nil{
			fmt.Println("error: ", err.Error(), n)
		}
		err = rw.Flush()
		if err != nil {
			fmt.Println("error: ", err.Error())
		}
		fmt.Println("send over!")
	}


	//for {
	//	var cData = datas.StructData{
	//		N:14,
	//		S:"wuxuntest",
	//	}
	//	jData, err := json.Marshal(cData)
	//	n, err := conn.Read(readBuffer)
	//	if err != nil {
	//		fmt.Println("Read failed:", err)
	//		return
	//	}
	//	_, err1 := conn.Write(jData)
	//	fmt.Println("count:", n, "msg:", readBuffer,  string(readBuffer), err1, writeBuffer)
	//
	//}

}
