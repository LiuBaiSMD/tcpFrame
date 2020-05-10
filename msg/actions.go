/*
@Author: liubai
@Date: 2020/5/10 6:16 下午
@Desc: use for what
*/

package msg

import (
	"bufio"
	"github.com/micro/go-micro/util/log"
	"tcpPractice/conns"
	"tcpPractice/datas"
	"tcpPractice/registry"
	"tcpPractice/util"
)

type ServerRfAddr struct {

}

func (b* ServerRfAddr)Communicate() registry.HttpWR{
	return func(rw *bufio.ReadWriter, BData datas.BaseData)error{
		log.Log("method:", util.RunFuncName()) //获取请求的方法
		return nil
	}
}

func (b* ServerRfAddr)HeartBeat() registry.HttpWR {
	return  func(rw *bufio.ReadWriter, BData datas.BaseData)error{
		log.Log("method:", util.RunFuncName()) //获取请求的方法
		SendMessage(rw, BData.Action, BData)
		if BData.UserId>0{
			conns.FlushConnLive(BData.UserId)
		}
		return nil
	}
}
