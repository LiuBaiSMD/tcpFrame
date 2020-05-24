/*
@Author: liubai
@Date: 2020/5/22 11:03 下午
@Desc: use for what
*/
package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"sort"
	"strconv"
	"strings"
)

func main() {
	var lastIndex uint64
	config := api.DefaultConfig()
	config.Address = "127.0.0.1:8500" //consul server

	client, err := api.NewClient(config)
	if err != nil {
		fmt.Println("api new client is failed, err:", err)
		return
	}
	serviceName := "serverNode"
	// 查看注册的服务信息
	services, metainfo, err := client.Health().Service(serviceName, "v3", true, nil)
	if err != nil {
		fmt.Println("error retrieving instances from Consul: %v", err)
		return
	}
	lastIndex = metainfo.LastIndex
	fmt.Println("lastIndex:", lastIndex)
	idList := make([]int, 0)
	for _, service := range services {
		fmt.Println("service.Service.Address:", service.Service)
		sId := service.Service.ID
		if strings.HasPrefix(sId, serviceName) {
			i, err := strconv.Atoi(sId[len(serviceName)+1:])
			if err != nil {
				continue
			}

			idList = append(idList, i)
		}
	}
	var lastId int
	sort.Ints(idList)
	if len(idList)<=0{
		lastId = 1
	}else{
		lastId = idList[len(idList)-1]+1
	}
	fmt.Println("idList: ", idList, lastId)
	client.Agent().ServiceDeregister("serverNode_")
}
