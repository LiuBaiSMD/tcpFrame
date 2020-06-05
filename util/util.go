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
	"io/ioutil"
	"net"
	"net/http"
	"runtime"
	"log"
	"fmt"
	"errors"
	mJson "github.com/micro/go-micro/config/encoder/json"
	"encoding/json"
)

func RunFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return "["+f.Name()+"]:"
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

func PanicErr(err error) bool {
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func LogErr(err error) bool {
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func GetBody(r *http.Request) (map[string]interface{}, error) {
	//解析get方法参数
	if r.Method == "GET" {
		r.ParseForm()
		res := make(map[string]interface{})
		for k, v := range r.Form {
			if len(v) != 1 {
				return nil, errors.New("请求参数重复，请检查！")
			}
			res[k] = v[0]
		}
		return res, nil
	}

	//将参数解析为 map[string]interface{}型
	if r.Method != "POST" {
		return nil, errors.New("请求类型错误，请检查")
	}
	ContType := r.Header["Content-Type"]
	if ContType[0] == "application/json" {
		if err := r.ParseForm(); err != nil {
			return nil, errors.New("参数解析异常")
		}
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, errors.New("连接错误")
		}
		var webData interface{}
		if err := json.Unmarshal(b, &webData); err != nil {
			return nil, errors.New("json解析异常")
		}
		mapdata := webData.(map[string]interface{})
		return mapdata, nil
	}
	if ContType[0] == "application/x-www-form-urlencoded" {
		r.ParseForm()
		var mapdata map[string]interface{}
		mapdata = make(map[string]interface{})
		for k, v := range r.Form {
			mapdata[k] = v[0]
		}
		return mapdata, nil
	}
	return nil, errors.New("请求HEADER类型错误，请检查！")
}

func CheckOKs(oks ...bool) bool {
	//检查oks是否全为true
	for _, v := range oks {
		if !v {
			return false
		}
	}
	return true
}
