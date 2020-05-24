/*
auth:   wuxun
date:   2020-05-07 20:07
mail:   lbwuxun@qq.com
desc:   how to use or use for what
*/

package util

import (
	"bytes"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source"
	"github.com/micro/go-micro/config/source/file"
	"net"
	"runtime"
	"log"
	"fmt"
	"errors"
	mJson "github.com/micro/go-micro/config/encoder/json"
)

func RunFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

func BytesCombine(pBytes ...[]byte) []byte {
	len := len(pBytes)
	s := make([][]byte, len)
	for index := 0; index < len; index++ {
		s[index] = pBytes[index]
	}
	sep := []byte("")
	return bytes.Join(s, sep)
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s:%s", msg, err)
		panic(fmt.Sprintf("%s:%s", msg, err))
	}
}

//指定文件中的配置
func ReadConfig(filePath string) (map[string]interface{}, error) {
	configPath := filePath
	e := mJson.NewEncoder()
	fileSource := file.NewSource(
		file.WithPath(configPath),
		source.WithEncoder(e),
	)
	conf := config.NewConfig()
	// 加载micro.yml文件
	if err := conf.Load(fileSource); err != nil {
		panic(err)
	}
	routes := make(map[string]interface{})
	err := conf.Scan(&routes)
	if err != nil {
		return nil, err
	}
	return routes, nil
}

//获取字典中的值
func GetMapContent(m map[string]interface{}, path ...string) (interface{}, error) {
	//本接口将获取一个map中，按path路径取值，返回一个interface
	var content interface{}
	var ok bool
	l := len(path)
	if l == 0 || (l == 1 && path[0] == "") { //当没有填入
		return m, nil
	}
	for k, v := range path {
		if k == l-1 {
			content, ok = m[v]
			if !ok {
				return nil, errors.New(" 配置读取错误---> 	" + v)
			}
			return content, nil
		}
		if m, ok = m[v].(map[string]interface{}); !ok {
			return nil, errors.New(" 配置读取错误---> 	" + v)
		}
	}
	return nil, errors.New("missing map!")
}

func localIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
