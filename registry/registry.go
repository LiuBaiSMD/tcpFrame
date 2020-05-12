// @Author: liubai
// @Date: 2020/5/10 11:28 上午
// @Desc: 自动注册处理方法，可以通过方法名字获取方法，并直接调用

package registry

//将传入的interf中的方法包装成map[string] HttpWR ，供路由绑定使用

import (
	"bufio"
	"fmt"
	"reflect"
	"tcpPractice/datas"
)
type ControllerMapsType map[string]reflect.Value
type HttpWR  func(w *bufio.ReadWriter,BData datas.BaseData) error
type Base struct{
	CrMap ControllerMapsType
	FuncRegistry map[string] HttpWR
}

var register *Base

func init(){
	register = &Base{}
	register.FuncRegistry = make(map[string] HttpWR)
	register.CrMap = make(ControllerMapsType, 0)
}

//注册函数，通过反射将handles中的方法，打包成字典存入 Register.FuncRegistry中 key为对应的方法名，value为对应的方法
func Registery(handles interface{})*Base{

	//创建反射变量，注意这里需要传入ruTest变量的地址；
	//不传入地址就只能反射Routers静态定义的方法
	vf := reflect.ValueOf(handles)
	vft := vf.Type()
	//读取方法数量
	mNum := vf.NumMethod()
	fmt.Println("NumMethod:", mNum)
	//遍历路由器的方法，并将其存入控制器映射变量中
	for i := 0; i < mNum; i++ {
		mName := vft.Method(i).Name
		register.CrMap[mName] = vf.Method(i)
		f:= register.CrMap[mName].Call(nil)
		_, ifOK := register.FuncRegistry[mName]
		if ifOK {
			//panic("重复注册方法 -----> " + mName)
		}
		register.FuncRegistry[mName] = f[0].Interface().(HttpWR)
	}
	if len(register.FuncRegistry) == 0{
		return nil
	}
	fmt.Println("FuncRegistry: ---->", register.FuncRegistry)
	return register
}

func GetHandleByName(funcName string)HttpWR{
	f, ok := register.FuncRegistry[funcName]
	if ok{
		return f
	}
	return nil
}

