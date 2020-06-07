/*
auth:   wuxun
date:   2019-12-09 19:54
mail:   lbwuxun@qq.com
desc:   用户的连接以及用户请求参数
*/

package conns

import (
	"bufio"
	"net"
)

type ClientConn struct {
	userId int `"用户id"`
	connID int `本次处理连接的id`
	conn   net.Conn
	rw     *bufio.ReadWriter
}

func NewClient(uId int, con net.Conn, cId int) *ClientConn {
	return &ClientConn{
		userId: uId,
		conn:   con,
		connID: cId,
		rw:     bufio.NewReadWriter(bufio.NewReader(con), bufio.NewWriter(con)),
	}
}

func (c ClientConn) GetUserID() int {
	return c.userId
}

func (c ClientConn) GetConnID() int {
	return c.connID
}

func (c ClientConn) GetConn() net.Conn {
	return c.conn
}

func (c ClientConn) GetRwBuf() *bufio.ReadWriter {
	return c.rw
}
