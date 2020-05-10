//此包用于实现具体web处理方法

package registry

import (
	"bufio"
	"github.com/micro/go-micro/util/log"
	"tcpPractice/datas"
	"tcpPractice/util"
)

type RfAddr struct {}

func Init(){
	var rfaddr RfAddr
	Registery(&rfaddr)
}

func (b* RfAddr)TestUserLogin() HttpWR{
	return func(w bufio.ReadWriter, BData datas.BaseData)error{
		log.Log("method:", util.RunFuncName()) //获取请求的方法
		return nil
	}
}

func (b* RfAddr) Login() HttpWR {
	return  func(w bufio.ReadWriter, BData datas.BaseData)error{
		log.Log("method:", util.RunFuncName()) //获取请求的方法
		return nil
	}
}

