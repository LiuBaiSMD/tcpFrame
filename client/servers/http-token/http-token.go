/*
@Author: liubai
@Date: 2020/6/3 11:22 下午
@Desc: 一个简单的http服务获取token
*/

package main

import (
	"fmt"
	"log"
	"net/http"
	"tcpFrame/config/consul"
	_const "tcpFrame/const"
	"tcpFrame/dao"
	"tcpFrame/handle"
	"tcpFrame/util"
)

var count int64

//本服务注册使用的ip和端口
var ipAddr string = "127.0.0.1"
var port int = 8081

func getTokenServer() {
	http.HandleFunc("/getToken", getToken)
	http.ListenAndServe(fmt.Sprintf("%s:%d", ipAddr, port), nil)

}

func getToken(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		r.ParseForm()
		fmt.Println(r.Body, r.Form)
		userId1, ok4 := r.Form["userId"]
		userName1, ok5 := r.Form["userName"]
		log.Println(userId1, userName1, ok4, ok5)

	}

	mapdata, err := util.GetBody(r)
	userId, ok1 := mapdata["userId"].(string)
	userName, ok2 := mapdata["userName"].(string)
	if err != nil || !util.CheckOKs(ok1, ok2) {
		log.Println("GetBody参数解析错误！", mapdata)
		fmt.Fprint(w, err)
		return
	}
	fmt.Println(util.RunFuncName(), userId, userName)
	// 首先获取用户token
	token, err := dao.GetuserToken(userId)
	if token != "" && err == nil {
		log.Println("just get")
		fmt.Fprint(w, token)
		return
	}

	token, err = handle.GetTokenReal(userId, userName)
	if err != nil {
		return
	}
	dao.SaveUserToken(userId, token)
	fmt.Fprint(w, token)
}

func main() {
	redisCfg, _ := consul.GetRedisCfg(_const.CONSUL_IP, _const.CONSUL_PORT)
	fmt.Println(util.RunFuncName(), redisCfg)
	if redisCfg == nil {
		panic("配置错误！")
	}
	fmt.Println(redisCfg)
	dao.InitRedis(redisCfg.Password, fmt.Sprintf("%s:%d", redisCfg.Ip, redisCfg.Port), redisCfg.DB)
	getTokenServer()
}
