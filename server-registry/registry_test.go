/*
@Author: liubai
@Date: 2020/5/24 7:34 下午
@Desc: use for what
*/

package server_registry_test

import (
	"fmt"
	sr "tcpFrame/server-registry"
	"tcpFrame/util"
	"testing"
)

func Test_Registry(t *testing.T) {
	//服务注册需要先连接上consul
	sr.ConsulConnect("localhost:8500")
	//注册一个服务 其ip、port、serviceName、tag信息
	sr.RegisterServer(
		"1.0.0.1",
		1234,
		"serverNode",
		[]string{},
		)
	//获取该serverName下的所有服务节点信息
	servicesMap, _ := sr.ServicesMap("serverNode")
	fmt.Println(util.RunFuncName(), servicesMap)
	//注销一个serverId为serverNode_3的节点
	sr.DeRegistry("serverNode_3")
	//注销serverName下的所有服务Id
	sr.DeRegistryAll("serverNode")
}
