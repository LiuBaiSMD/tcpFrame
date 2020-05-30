package consul_test

import (
	"fmt"
	configCs "tcpFrame/config/consul"
	"tcpFrame/util"
	"testing"
)

func Test_GetConsulConfig(t *testing.T) {
	config, err := configCs.ReaderConfig("127.0.0.1", 8500, []string{"testConfig", "consul_config"})
	fmt.Println(util.RunFuncName(), config, err)
	fmt.Println(util.RunFuncName(), string(config))
}

func Test_PutConsulConfig(t *testing.T) {
	err := configCs.UpdataConfig("127.0.0.1", 8500, ".", "service.json", "serverRegistry")
	fmt.Println(util.RunFuncName(), err)
}
