// @Author: liubai
// @Date: 2020/5/10 6:11 下午
// @Desc: use for what

package msg

import (
	"bufio"
	"fmt"
	"net"
	"tcpFrame/datas"
	"tcpFrame/registry"
	"tcpFrame/util"
)

var register *registry.Base

func init(){
	var rfaddr1 ServerRfAddr
	register = registry.Registery(&rfaddr1)
}

//根据协议中的action进行分包, 早出对应的funcname进行处理
func DisPatch(conn net.Conn, data interface{}) error{
	cData, ok := data.(datas.BaseData)
	fmt.Println(util.RunFuncName(), string(cData.BData))
	funcName := cData.Action
	if funcName==""{
		fmt.Println("action is empty")
	}
	f := registry.GetHandleByName(funcName)
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	if f==nil{
		fmt.Printf("func %s not in registry ", funcName)
		SendMessage(rw, cData.Action, cData)
		return nil
	}
	f(rw, cData)
	fmt.Println("loginAfter", cData, ok)
	return nil
}


