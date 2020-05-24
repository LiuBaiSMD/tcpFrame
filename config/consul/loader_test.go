package consul_test

import (
	"fmt"
	configCs "tcpFrame/config/consul"
	"tcpFrame/util"
	"testing"
)

func Test_GetConsulConfig(t *testing.T) {
	config, err := configCs.ReaderConfig("127.0.0.1", 8500, []string{"serviceConfig", "consul_config"})
	fmt.Println(util.RunFuncName(), config, err)
	fmt.Println(util.RunFuncName(), string(config))
}

func Test_PutConsulConfig(t *testing.T) {
	err := configCs.UpdataConfig("127.0.0.1", 8500, ".", "totalConfig.json", "testConfig")
	fmt.Println(util.RunFuncName(), err)
}

type consulConfig struct{
	Enabled     bool    `json:"enabled"`
	Host 		string  `josn:"host"`
	Port 		int		`json:"port"`
	DockerHost  string	`json:"docker_host"`
}