/*
auth:   wuxun
date:   2019-12-09 17:25
mail:   lbwuxun@qq.com
desc:   1.store conns (push and pop)
		2.return the current connID which is pop just now
		3.return the current lenth of channel of conns
*/

package conns

import (
	"fmt"
	"sync"
	"tcpFrame/dao"
	"tcpFrame/util"
	"time"
)

type ConnMap struct {
	//connChan chan int
	connMap  sync.Map
	curConnID int
	connCMap chan *ClientConn
	connLiveMap sync.Map
}

var connIDCreator chan int
var cMap ConnMap
var connCount int

func init() {
	cMap.connCMap = make(chan *ClientConn, 10000)
	connIDCreator = make(chan int, 1)
	cMap.curConnID = -1
	connIDCreator <- 1
	dao.Init()
	//关闭心跳检测
	go manageConnLive()
}

func PushChan(connID int, connValue interface{}){
	// 如果已经有了连接，先断开此链接
	oldConnCli := GetConnByUId(int(connID))
	if oldConnCli!=nil{
		oldConn := oldConnCli.GetConn()
		oldConn.Close()
	}
	cMap.connMap.Store(connID, connValue)
	cMap.connLiveMap.Store(connID, time.Now().Unix())
}

func FlushConnLive(connID int){
	fmt.Println(util.RunFuncName(), connID, time.Now().Unix())
	conn := GetConnByUId(connID)
	if conn==nil{
		return
	}
	cMap.connLiveMap.Store(connID, time.Now().Unix())
}

func GetConnByUId(connId int)*ClientConn{
	if connId<1{
		return nil
	}
	connValue, isOK := cMap.connMap.Load(connId)
	if !isOK{
		return nil
	} else {
		return connValue.(*ClientConn)
	}
}

func DelConnById(cId int){
	//先断开连接
	if cli, isOk := cMap.connMap.Load(cId);isOk{
		connCli := cli.(*ClientConn)
		connCli.conn.Close()
	}

	cMap.connMap.Delete(cId)
	cMap.connLiveMap.Delete(cId)

}

func LenthConn()int{
	connCount = 0
	cMap.connLiveMap.Range(countLenth)
	//return len()
	return connCount
}

func countLenth(k, v interface{})bool{
	connCount++
	return true
}

func GetCMap()ConnMap{
	return cMap
}