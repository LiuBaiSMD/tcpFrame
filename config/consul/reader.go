/*
auth:   wuxun
date:   2020-05-23 18:55
mail:   lbwuxun@qq.com
desc:   how to use or use for what
*/

package consul

import (
	"fmt"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source/consul"
	"github.com/micro/go-micro/util/log"
)

//func ReadConfig

//通过存储配置的consul服务ip 端口、以及路径获取consul配置，返回[]byte供使用者自行序列化
func ReaderConfig(consulServerIp string, port int, configAddrList []string) ([]byte, error) {
	consulSource := consul.NewSource(
		consul.WithAddress(fmt.Sprintf("%s:%d", consulServerIp, port)),
		consul.WithPrefix(""),
	)
	// 创建新的配置
	conf := config.NewConfig()
	if err := conf.Load(consulSource); err!=nil {
		log.Logf("load config errr!", err)
		return nil, err
	}

	confData := conf.Get(configAddrList...).Bytes()
	return confData, nil
}
