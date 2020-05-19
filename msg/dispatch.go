// @Author: liubai
// @Date: 2020/5/10 6:11 下午
// @Desc: use for what

package msg

import (
	"bufio"
	"fmt"
	"net"
	"tcpFrame/const"
	"tcpFrame/datas"
	"tcpFrame/datas/proto"
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
		rsp := &heartbeat.LoginRespone{
			Code:200,
			LoginState:1,
			Oms:"login success!",
		}
		SendMessage(rw, _const.CMD_LOGIN_RSP, _const.BT_LOGIN_RSP, rsp)
		return nil
	}
	f(rw, cData)
	fmt.Println("loginAfter", cData, ok)
	return nil
}


