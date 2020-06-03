/*
auth:   wuxun
date:   2020-05-11 15:50
mail:   lbwuxun@qq.com
desc:   how to use or use for what
*/

package msg

import (
	"encoding/binary"
	"github.com/gogo/protobuf/proto"
	"tcpFrame/const"
	"tcpFrame/datas/proto"
	"tcpFrame/util"
)

var IoBuf []byte

/*
包的解析：
tcp中传输的二进制数据可以概括为
头部长度（四个字节）+ 消息长度（四个字节）+ 头部二进制数据 + 消息二进制数据
头部二进制数据序列化为 MsgBody结构体，ServerType：其中记录接受的服务名字  CmdType:序列化消息的类型（也可代表动作） UserId:发送的用户Id
*/

//从本地[]byte缓存解析出完整的原始[]byte数据包 headerLen, msgLen
func Parse2HeaderAndMsg(ioBuf *[]byte) (headerBytes, msgBytes []byte, err error) {
	//fmt.Println("Parse2HeaderAndMsg byte: ", ioBuf)
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
	//fmt.Println("tbData: ", tbData)
	return tbData, nil
}

func DecodeLenth(i uint32) []byte {
	var buf = make([]byte, _const.LEN_INFO_BYTE_SIZE)
	binary.BigEndian.PutUint32(buf, i)
	return buf
}

func EncodeLenthByte(buf []byte) uint32 {
	return binary.BigEndian.Uint32(buf)
}

func ParseMsg2RbtByte(senderId, cmdType string, userId int64, msgType int32, msgBytes []byte) []byte {
	rbtByte := &heartbeat.MsgBody{
		MsgType: msgType,
		SenderId: senderId,
		CmdType: cmdType,
		UserId: userId,
		MsgBytes: msgBytes,
		Version: version,
	}
	bData, _ := proto.Marshal(rbtByte)
	return bData
}
