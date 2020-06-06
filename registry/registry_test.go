// @Author: liubai
// @Date: 2020/5/10 11:54 上午
// @Desc: use for what

package registry_test

import (
	"bufio"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/micro/go-micro/util/log"
	heartbeat "tcpFrame/datas/proto"
	"tcpFrame/registry"
	"tcpFrame/util"
	"testing"
)

type RfAddrTest struct {

}

func (b* RfAddrTest)TestUserLogintest() registry.HttpWR{
	return func(w *bufio.ReadWriter, BData []byte)error{
		log.Log("method:", util.RunFuncName()) //获取请求的方法
		return nil
	}
}

func (b* RfAddrTest) Logintest() registry.HttpWR {
	return  func(w *bufio.ReadWriter, BData []byte)error{
		log.Log("method:", util.RunFuncName()) //获取请求的方法
		return nil
	}
}

func Test_Registry(t *testing.T) {
	var rfaddr RfAddrTest
	register := registry.Registery(&rfaddr)
	if register!= nil{
		fmt.Println(util.RunFuncName(), register.FuncRegistry["Logintest"])
		f := register.FuncRegistry["Logintest"]
		data := &heartbeat.LoginRequest{}
		bData, _ := proto.Marshal(data)
		f(&bufio.ReadWriter{}, bData)
	}
	var rfaddr1 RfAddrTest
	register1 := registry.Registery(&rfaddr1)
	if register1!= nil{
		fmt.Println(util.RunFuncName(), register1.FuncRegistry["Logintest"])
		f := register1.FuncRegistry["Logintest"]
		data := &heartbeat.LoginRequest{}
		bData, _ := proto.Marshal(data)
		f(&bufio.ReadWriter{}, bData)
	}
	for key, _ :=range register.FuncRegistry{
		fmt.Println("func key: ", key)
	}
	for key, _ :=range register1.FuncRegistry{
		fmt.Println("func key: ", key)
	}
}

