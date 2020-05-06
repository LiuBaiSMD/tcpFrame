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
	"tcpPractice/datas"
	"tcpPractice/msg"
	"time"
)

func Open(addr string) (*bufio.ReadWriter, net.Conn, error) {
	// Dial the remote process.
	// Note that the local port is chosen on the fly. If the local port
	// must be a specific one, use DialTCP() instead.
	fmt.Println("Dial " + addr)
	//conn, err := tls.Dial("tcp", addr, nil)
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
	var cData = datas.Request{
		Action:"login",
		Name:"wuxun",
		PWD:"123456",
		UserId:10001,
	}

	go msg.ListenMessageClient(conn)

	for{
		bData, _ := json.Marshal(cData)
		_, err := conn.Write(bData)
		if err!=nil{
			conn.Close()
			return
		}
		time.Sleep(time.Second * 5)
		fmt.Println("send over!")
	}

}
