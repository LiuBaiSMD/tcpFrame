/*
auth:   wuxun
date:   2020-05-11 15:50
mail:   lbwuxun@qq.com
desc:   how to use or use for what
*/

package msg

import (
	"encoding/binary"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"tcpFrame/const"
	heartbeat "tcpFrame/datas/proto"
	"tcpFrame/util"
)

var IoBuf []byte

/*
## 6.改进数据包传输协议
### ①增加组装传输数据的接口
```
总共分为两层
1.(第一层解析)数据包长度dataLenth（32位 []byte）+ 编码类型codeType（8位 []byte）+ 数据data（[]byte）
dataLenth:存储data长度
codeType:基础解析格式，标识解析data的方式，json、proto等通用的格式
data:数据内容

2.(第二层解析)解析data模块，将data分解成各个类型json、proto等的BaseData后，其中的Action数据为指导业务层自行解析的模块，比如
例① json中的BaseData结构:
type BaseData struct{
    Action string,
    UserId int,
    BData []byte,
}

json中的HeartBeat结构:
type HeartBeat struct{
    Action string,
    UserId int,
    TimeStamp int,
    OtherMsg string,
}

例如在上述拆包过程中codeType=1代表json格式数据，将data解析为json的BaseData格式:得到以下数据
data = {
        Action:"Heartbeat",
        UserId:10001,
        BData:[12, 23, 45, 234, 54, 65],
        }
(第三层解析)然后业务层通过Action将指导BData解析为已经定义好的json结构 HeartBeat
BData = {
    Action: HeartBeat,
    UserId: 10001,
    TimeStamp: 123456789,
    OtherMsg: "hello world!",
}

例②
proto中的BaseData结构:
message BaseData {
    string Action = 1;
    int64 UserId = 2;
    bytes BData = 3;
}

proto中的HeartBeat结构:
message HeartBeat {
    string Action = 1;
    int64 UserId = 2;
    int64 TimeStamp = 3;
    string OtherMsg = 4;
}

例如在上述拆包过程中codeType=2代表proto格式数据，将data解析为proto的BaseData格式:得到以下数据
data = {
        Action:"Heartbeat",
        UserId:10001,
        BData:[12, 23, 45, 234, 54, 65],
        }
(第三层解析)然后业务层通过Action将BData解析为已经定义好的proto结构 HeartBeat，
BData = {
    Action: HeartBeat,
    UserId: 10001,
    TimeStamp: 123456789,
    OtherMsg: "hello world!",
}
*/

//从本地[]byte缓存解析出完整的原始[]byte数据包 headerLen, msgLen
func Parse2HeaderAndMsg(ioBuf *[]byte) (headerBytes, msgBytes []byte, err error) {
	fmt.Println("Parse2HeaderAndMsg byte: ", ioBuf)
	//使用for循环模拟一次完整的数据读取
	if len(*ioBuf) <= (2 * _const.LEN_INFO_BYTE_SIZE) {
		return []byte(""), []byte(""), nil
	}
	offset := 0
	//先读取一个lenthData 4个字节
	headerLen := int(EncodeLenthByte((*ioBuf)[offset : offset+_const.LEN_INFO_BYTE_SIZE]))
	//读取编码格式codeType 4个字节
	offset = offset + _const.LEN_INFO_BYTE_SIZE
	msgLen := int(EncodeLenthByte((*ioBuf)[offset : offset+_const.LEN_INFO_BYTE_SIZE]))

	offset = offset + _const.LEN_INFO_BYTE_SIZE
	if len(*ioBuf) < (offset + headerLen + msgLen) {
		return []byte(""), []byte(""), nil
	}
	headerBytes = (*ioBuf)[offset : offset+headerLen]
	offset = offset + headerLen
	msgBytes = (*ioBuf)[offset : offset+msgLen]

	*ioBuf = (*ioBuf)[offset+msgLen:]
	return headerBytes, msgBytes, nil
}

//组装一个header长度和body长度到头部中去，根据codetype进行编码marshal rawdata
func BuildData(headerBytes, msgBytes []byte) ([]byte, error) {
	//序列化rawData

	//量出headerBytes长度，再将长度进行序列化，长度为4的[]byte
	headerLenInfo := (uint32)(len(headerBytes))
	headerLenInfoBytes := DecodeLenth(headerLenInfo)

	//量出rawData长度，再将长度进行序列化，长度为4的[]byte
	msgLenInfo := (uint32)(len(msgBytes))
	msgLenInfoBytes := DecodeLenth(msgLenInfo)

	//组装headerlen bodylen header msg信息
	tbData := util.BytesCombine(headerLenInfoBytes, msgLenInfoBytes, headerBytes, msgBytes)
	fmt.Println("tbData: ", tbData)
	return tbData, nil
}

func DecodeLenth(i uint32) []byte {
	var buf = make([]byte, _const.LEN_INFO_BYTE_SIZE)
	binary.BigEndian.PutUint32(buf, i)
	return buf
}

func EncodeLenthByte(buf []byte) uint32 {
	return uint32(binary.BigEndian.Uint32(buf))
}

func ParstMsg2RbtByte(cmdType string, msgBytes []byte) []byte{
	rbtByte := &heartbeat.MsgBody{
		CmdType: cmdType,
		MsgBytes: msgBytes,
	}
	bData, _ := proto.Marshal(rbtByte)
	return bData
}
