/*
auth:   wuxun
date:   2020-05-11 15:55
mail:   lbwuxun@qq.com
desc:   how to use or use for what
*/

package msg_test

import (
	"fmt"
	"tcpFrame/util"
	"testing"
	"tcpFrame/msg"
	"tcpFrame/datas"
)

var ioBuf []byte

func Test_ExchangeData(t *testing.T) {
	ioBuf, _ = msg.BuildData(1, datas.BaseData{Action:"test", UserId:10001, BData:[]byte("Hello world!")})
	ioBuf1, _ := msg.BuildData(1, datas.BaseData{Action:"tsetBuildData", UserId:10001})
	for i:=0;i<len(ioBuf1);i++{
		ioBuf = append(ioBuf, ioBuf1[i])
	}
	msg.IoBuf = ioBuf
	fmt.Println(util.RunFuncName(), ioBuf)
	codeType, bRawData, err :=msg.ReadData(&ioBuf)
	fmt.Println(codeType, bRawData, err)
	codeType, bRawData, err =msg.ReadData(&ioBuf)
	fmt.Println(codeType, bRawData, err)
}