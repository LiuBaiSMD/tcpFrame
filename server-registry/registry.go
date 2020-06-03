/*
@Author: liubai
@Date: 2020/5/24 5:59 下午
@Desc: use for what
*/

package server_registry

import (
	"errors"
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"log"
	_ "net/http/pprof"
	"sort"
	"strconv"
	"strings"
	"tcpFrame/util"
	"os"
	"os/signal"
	"syscall"
)

var csCli *consulapi.Client
var servers map[string]*consulapi.AgentServiceRegistration

func init() {
	servers = make(map[string]*consulapi.AgentServiceRegistration)
}

func ConsulConnect(consulUrl string) {
	config := consulapi.DefaultConfig()
	config.Address = consulUrl
	csCli, _ = consulapi.NewClient(config)
}

func DeferDeregistry(serverId string) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c //阻塞等待
		DeRegistry(serverId)
		os.Exit(0)
	}()
}

func RegisterServer(serviceIp string, servicePort int, serviceName string, tags []string) ( string, error) {
	if csCli == nil {
		print("consul client not init!")
		return "", errors.New("consul client not init!")
	}
	//注册前确定service的Name，以及在此Name下的serviceId
	registration := &consulapi.AgentServiceRegistration{}
	serviceId, _ := getLastServiceIdByServiceName(csCli, serviceName)
	registration.ID = serviceId      // 服务节点的名称
	registration.Name = serviceName  // 服务名称
	registration.Port = servicePort  // 服务端口
	registration.Tags = tags         // tag，可以为空
	registration.Address = serviceIp // 服务 IP
	registration.Kind = "getToken"   //kind,服务类型
	err := csCli.Agent().ServiceRegister(registration)
	if err != nil {
		log.Fatal("register server error : ", err)
		return "", err
	}
	servers[registration.ID] = registration
	fmt.Println(util.RunFuncName(), servers)
	go DeferDeregistry(serviceId)
	return serviceId, nil
}

func getLastServiceIdByServiceName(cs *consulapi.Client, serviceName string) (string, error) {
	services, _, err := cs.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return "", err
	}
	idList := make([]int, 0)
	for _, service := range services {
		sId := service.Service.ID
		if strings.HasPrefix(sId, serviceName) {
			i, err := strconv.Atoi(sId[len(serviceName)+1:])
			if err != nil {
				continue
			}
			idList = append(idList, i)
		}
	}
	sort.Ints(idList)
	var lastId int
	if len(idList) <= 0 {
		lastId = 1
	} else {
		lastId = idList[len(idList)-1] + 1
	}
	serviceNodeName := fmt.Sprintf("%s_%d", serviceName, lastId)
	fmt.Println(util.RunFuncName(), serviceNodeName, lastId)
	return serviceNodeName, nil
}

func DeRegistry(serverId string) error {
	err := csCli.Agent().ServiceDeregister(serverId)
	for _, ok := servers[serverId]; !ok; {
		delete(servers, serverId)
		break
	}
	return err
}

func DeRegistryAll(serverName string) error {
	services, _, err := csCli.Health().Service(serverName, "", true, nil)
	if err != nil {
		return err
	}
	for _, service := range services {
		err = DeRegistry(service.Service.ID)
		if err != nil {
			return nil
		}
	}
	return nil
}

func ServicesMap(serverName string) (map[string]*consulapi.ServiceEntry, error) {
	services, _, err := csCli.Health().Service(serverName, "", true, nil)
	if err!=nil{
		return nil, err
	}
	serviceMap := make(map[string]*consulapi.ServiceEntry)
	for _, service := range services {
		fmt.Println("service.Service.Address:", service.Service)
		sId := service.Service.ID
		serviceMap[sId] = service
	}
	return serviceMap, err
}