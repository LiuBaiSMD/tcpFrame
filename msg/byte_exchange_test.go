/*
auth:   wuxun
date:   2020-05-11 15:55
mail:   lbwuxun@qq.com
desc:   how to use or use for what
*/

package msg_test
import (
	"testing"
	"tcpPractice/msg"
	"tcpPractice/datas"
)
var ioBuf []byte

func Test_ExchangeData(t *testing.T) {
	ioBuf, _ = msg.BuildData(1, datas.BaseData{Action:"test", UserId:10001})
	ioBuf1, _ := msg.BuildData(1, datas.BaseData{Action:"tsetBuildData", UserId:10001})
	for i:=0;i<len(ioBuf1);i++{
		ioBuf = append(ioBuf, ioBuf1[i])
	}
	msg.ReadData()
	msg.ReadData()
}