/*
auth:   wuxun
date:   2020-05-11 15:50
mail:   lbwuxun@qq.com
desc:   how to use or use for what
*/

package msg

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"tcpPractice/util"
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

//从本地[]byte缓存解析出完整的原始[]byte数据包
func ReadData(ioBuf *[]byte)(codeType int,bRawData []byte,err error){
	fmt.Println("readData byte: ", ioBuf)
	//使用for循环模拟一次完整的数据读取
	if len(*ioBuf)<=5{
		return 0, []byte(""), nil
	}
	//先读取一个lenthData 4个字节
	lenthData := (*ioBuf)[0:4]
	//读取编码格式codeType 1个字节
	codeType = int((*ioBuf)[4])
	//根据lenthData 读取对应唱的的data
	l := LenthToInt(lenthData)
	//根据codeType 解析数据，如果数据长度不够，则不作处理，返回空值不报错
	if len(*ioBuf)<5+l{
		fmt.Println("l: ", len(*ioBuf), l)
		return 0, []byte(""), nil
	}
	fmt.Println("l: ", len(*ioBuf), l)
	bRawData = (*ioBuf)[5:5+l]
	*ioBuf = (*ioBuf)[5+l:]
	return int(codeType), bRawData, nil
}

//组装一个长度和编码类型到头部中去，根据codetype进行编码marshal rawdata
func BuildData(codeType int,rawData interface{})([]byte, error){
	//序列化rawData
	bRawData, _ := DecodeData(rawData)

	//量出rawData长度，再将长度进行序列化，长度为4的[]byte
	l := (int32)(len(bRawData))
	lRawData := DecodeLenth(l)

	//将序列化类型codeType序列化，长度为1的[]byte
	cRawData := DecodeCodeType(int8(codeType))
	//组装一个数据包
	tbData := util.BytesCombine(lRawData, cRawData, bRawData)
	fmt.Println("tbData: ", tbData)
	return tbData, nil
}

func DecodeData(rawData interface{})([]byte, error){
	bData, err := json.Marshal(rawData)
	return bData, err
}

func DecodeLenth(i int32) []byte {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(i))
	return buf
}

func DecodeCodeType(i int8) []byte {
	var buf = make([]byte, 1)
	buf[0] = byte(i)
	return buf
}

func LenthToInt(buf []byte) int {
	return int(binary.BigEndian.Uint32(buf))
}
