// @Author: liubai
// @Date: 2020/5/2 5:26 下午
// @Desc: use for what

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

func main() {

	addr := "0.0.0.0:8080"

	tcpAddr, err := net.ResolveTCPAddr("tcp",addr)

	if err != nil {
		log.Fatalf("net.ResovleTCPAddr fail:%s", addr)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("listen %s fail: %s", addr, err)
	} else {

		log.Println("rpc listening", addr)
	}


	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("listener.Accept error:", err)
			continue
		}

		go handleConnection(conn)

	}

}

type ComplexData struct{
	N int
	S  string
}

func handleConnection(conn net.Conn) {
	//根据连接的数据进行dispach

	defer conn.Close()

	//readBuffer := make([]byte, 512)
	//var writeBuffer []byte = []byte("You are welcome. I'm server.")
	bData := make([]byte, 10)
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	for {
		bData1, err := rw.Read(bData)
		if err != nil{
			fmt.Println("链接无法读取.", err)
			return
		}
		if bData1>0{
			var cData ComplexData
			err:=json.Unmarshal(bData[:bData1], &cData)
			if err != nil{
				fmt.Println("err:", err)
			}
			fmt.Println("respone: ", bData, bData1, string(bData), cData.N, cData.S)
		}
		//// 写入底层网络链接
		//err = rw.Flush()
		//if err != nil{
		//	fmt.Println("Flush写入失败")
		//	return
		//}
	}
}
