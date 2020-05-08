/*
auth:   wuxun
date:   2019-12-09 19:54
mail:   lbwuxun@qq.com
desc:   用户的连接以及用户请求参数
*/

package conns

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/micro/go-micro/util/log"
	"net"
	"tcpPractice/proto"
)

type ClientConn struct{
	userId int				`"用户id"`
	connID int				`本次处理连接的id`
	conn net.Conn
}
func NewClient(uId int, con net.Conn, cId int)  *ClientConn{
	return &ClientConn{
		userId:uId,
		conn: con,
		connID:cId,
	}
}

func (c ClientConn)GetUserID()int{
	return c.userId
}

func (c ClientConn)GetConnID()int{
	return c.connID
}

func (c ClientConn)GetConn()net.Conn{
	return c.conn
}
func (c ClientConn)ListenMessage(){
	done := make(chan struct{})
	readBuffer := make([]byte, 128)
	clientRes := heartbeat.Response{}
	go func() {
		defer close(done)
		for {
			n, err := c.conn.Read(readBuffer)
			if err != nil {
				log.Log("read:", err, n)
				return
			}
			if err := proto.Unmarshal(readBuffer, &clientRes); err != nil {
				log.Logf("proto unmarshal: %s", err)
			}
			fmt.Println("recv: ", clientRes.Data)
		}
	}()
}

