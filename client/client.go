// @Author: liubai
// @Date: 2020/5/2 5:27 下午
// @Desc: use for what

package main


import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type ComplexData struct{
	N int
	S  string
}

func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("dial failed:", err)
		os.Exit(1)
	}
	defer conn.Close()

	readBuffer := make([]byte, 512)
	writeBuffer := []byte("i am connector!")

	for {
		var cData = ComplexData{
			N:14,
			S:"wuxun",
		}
		jData, err := json.Marshal(cData)
		n, err := conn.Read(readBuffer)
		if err != nil {
			fmt.Println("Read failed:", err)
			return
		}
		_, err1 := conn.Write(jData)
		fmt.Println("count:", n, "msg:", string(readBuffer), err1, writeBuffer)

	}

}
