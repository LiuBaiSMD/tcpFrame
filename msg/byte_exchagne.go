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
	"errors"
	"fmt"
	"tcpPractice/util"
)
var IoBuf []byte
func ReadData(ioBuf *[]byte)error{
	fmt.Println("readData byte: ", ioBuf)
	//使用for循环模拟一次完整的数据读取
	for{
		if len(*ioBuf)<=5{
			break
		}
		//先读取一个lenthData 4个字节
		lenthData := (*ioBuf)[0:4]
		//读取编码格式codeType 1个字节
		codeType := (*ioBuf)[4]
		//根据lenthData 读取对应唱的的data
		l := LenthToInt(lenthData)
		//根据codeType 解析数据
		if len(*ioBuf)<5+l{
			fmt.Println("l: ", len(*ioBuf), l)
			continue
		}
		fmt.Println("l: ", len(*ioBuf), l)
		bRawData := (*ioBuf)[5:5+l]
		*ioBuf = (*ioBuf)[5+l:]
		fmt.Println("readData: ", lenthData, codeType, l, bRawData, *ioBuf)
		fmt.Println("get data: ", string(bRawData))
		break
	}
	return nil
}

//11001100010010110111111011111000011011111011011
//组装一个长度和编码类型到头部中去，根据codetype进行编码marshal rawdata
func BuildData(codeType int,rawData interface{})([]byte, error){
	//序列化rawData
	bRawData, _ := DecodeData(rawData)

	//量出rawData长度，再将长度进行序列化，长度为4的[]byte
	l := (int32)(len(bRawData))
	lRawData := DecodeLenth(l)

	//将序列化类型codeType序列化，长度为1的[]byte
	cRawData := DecodeCodeType(int8(codeType))

	fmt.Println("lenthData: ", lRawData)
	fmt.Println("codeType: ", cRawData)
	fmt.Println("rawData: ", bRawData)

	//组装一个数据包
	tbData := util.BytesCombine(lRawData, cRawData, bRawData)
	fmt.Println("tbData: ", tbData)
	return tbData, errors.New("test")
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
