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
	"tcpPractice/dao"
	"tcpPractice/util"
	"time"
)

type ConnMap struct {
	//connChan chan int
	connMap  sync.Map
	curConnID int
	connCMap chan *ClientConn
	connLiveMap map[int]int64
}

var connIDCreator chan int
var cMap ConnMap

func init() {
	cMap.connCMap = make(chan *ClientConn, 10000)
	connIDCreator = make(chan int, 1)
	cMap.connLiveMap = make(map[int]int64)
	cMap.curConnID = -1
	connIDCreator <- 1
	dao.Init()
	go manageConnLive()
}

func PushChan(connID int, connValue interface{}){
	cMap.connMap.Store(connID, connValue)
	cMap.connLiveMap[connID] = time.Now().Unix()
}

func FlushConnLive(connID int){
	fmt.Println(util.RunFuncName(), connID, time.Now().Unix())
	conn := GetConnByUId(connID)
	if conn==nil{
		return
	}
	cMap.connLiveMap[connID] = time.Now().Unix()
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
	conn := GetConnByUId(cId)
	conn.GetConn().Close()
	cMap.connMap.Delete(cId)
	delete(cMap.connLiveMap, cId)

}

func LenthConn()int{
	return len(cMap.connLiveMap)
}

func GetCMap()ConnMap{
	return cMap
}