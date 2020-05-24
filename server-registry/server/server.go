/*
@Author: liubai
@Date: 2020/5/22 11:03 下午
@Desc: use for what
*/

package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"

	consulapi "github.com/hashicorp/consul/api"
)

var count int64

// consul 服务端会自己发送请求，来进行健康检查
func consulCheck(w http.ResponseWriter, r *http.Request) {

	s := "consulCheck" + fmt.Sprint(count) + "remote:" + r.RemoteAddr + " " + r.URL.String()
	fmt.Println(s)
	fmt.Fprintln(w, s)
	count++
}

func registerServer() {

	config := consulapi.DefaultConfig()
	config.Address = "127.0.0.1:8500"
	client, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatal("consul client error : ", err)
	}

	registration := &consulapi.AgentServiceRegistration{}
	registration.ID = "serverNode_2"      // 服务节点的名称
	registration.Name = "serverNode"      // 服务类型名称
	registration.Port = 9527              // 服务端口
	registration.Tags = []string{"v2000"} // tag，可以为空
	registration.Address = localIP()      // 服务 IP

	registration2 := &consulapi.AgentServiceRegistration{}
	registration2.ID = "serverNode_3"      // 服务节点的名称
	registration2.Name = "serverNode"      // 服务类型名称
	registration2.Port = 9527              // 服务端口
	registration2.Tags = []string{"v2000", "v3"} // tag，存放版本号等信息
	registration2.Address = localIP()      // 服务 IP
	//checkPort := 8081 // 给注册的服务加上健康检查方法
	//registration.Check = &consulapi.AgentServiceCheck{ // 健康检查
	//	HTTP:                           fmt.Sprintf("http://%s:%d%s", registration.Address, checkPort, "/check"),
	//	//HTTP:                           "www.baidu.com",
	//	Timeout:                        "3s",
	//	Interval:                       "5s",  // 健康检查间隔
	//	DeregisterCriticalServiceAfter: "5s", //check失败后30秒删除本服务，注销时间，相当于过期时间
	//	// GRPC:     fmt.Sprintf("%v:%v/%v", IP, r.Port, r.Service),// grpc 支持，执行健康检查的地址，service 会传到 Health.Check 函数中
	//}

	err = client.Agent().ServiceRegister(registration)
	err = client.Agent().ServiceRegister(registration2)
	if err != nil {
		log.Fatal("register server error : ", err)
	}
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

func main() {
	registerServer()
}
