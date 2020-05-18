/*
@Author: liubai
@Date: 2020/5/18 10:02 下午
@Desc: 数据解析模块，解析基础请求头，解析普通数据
*/

package msg

import (
	"tcpFrame/datas/proto"
	"unsafe"
)

func RequestMinLen()int{
	rqsHeader := &heartbeat.RequestHeader{}
	return int(unsafe.Sizeof(rqsHeader.CmdNo)+unsafe.Sizeof(rqsHeader.HeadLength)+unsafe.Sizeof(rqsHeader.BodyLength)+unsafe.Sizeof(rqsHeader.BodyType)+unsafe.Sizeof(rqsHeader.Version))
}