// @Author: liubai
// @Date: 2020/5/2 5:26 下午
// @Desc: use for what

package main

import (
	"log"
	"net"
	"time"

	"github.com/felixge/tcpkeepalive"
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

func setTcpKeepAlive(conn net.Conn) (*tcpkeepalive.Conn, error) {

	newConn, err := tcpkeepalive.EnableKeepAlive(conn)
	if err != nil {
		log.Println("EnableKeepAlive failed:", err)
		return nil, err
	}

	err = newConn.SetKeepAliveIdle(10*time.Second)
	if err != nil {
		log.Println("SetKeepAliveIdle failed:", err)
		return nil, err
	}


	err = newConn.SetKeepAliveCount(9)
	if err != nil {
		log.Println("SetKeepAliveCount failed:", err)
		return nil, err
	}

	err = newConn.SetKeepAliveInterval(10*time.Second)
	if err != nil {
		log.Println("SetKeepAliveInterval failed:", err)
		return nil, err
	}

	return newConn, nil
}


func handleConnection(conn net.Conn) {
	defer conn.Close()

	newConn, err := setTcpKeepAlive(conn)
	if err != nil {
		log.Println("setTcpKeepAlive failed:", err)
		return
	}
	readBuffer := make([]byte, 512)
	var writeBuffer []byte = []byte("You are welcome. I'm server.")


	for {

		time.Sleep(1*time.Second)
		n, err := newConn.Write(writeBuffer)
		if err != nil {
			log.Println("Write error:", err)
			break
		}
		log.Println("send:", n)
		n1, err1 := newConn.Read(readBuffer)
		log.Println("get:", n1, err1, string(readBuffer))
		select{

		}
	}

	log.Println("connetion end")

}
