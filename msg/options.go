/*
@Author: liubai
@Date: 2020/6/7 3:15 下午
@Desc: use for what
*/

package msg

import (
	"encoding/json"
	"fmt"
	"log"
	configCs "tcpFrame/config/consul"
	"tcpFrame/const"
	"tcpFrame/dao"
	"tcpFrame/msgMQ/nats-mq"
	"tcpFrame/registry"
	"tcpFrame/util"
)

type opt func(o *opts)

type opts struct {
	redisCfg    configCs.RedisConfig
	consulCfg   configCs.ConsulConfig
	natsCfg     configCs.NatsConfig
	rabbitmqCfg configCs.RabbitmqConfig
}

var myOpts opts

//tcpConn连接服注册方法
func InitServer(serverId string, setOpts ...opt) {

	for _, o := range setOpts {
		o(&myOpts)
	}
	//初始化数据库
	dao.InitRedis(myOpts.redisCfg.Password, fmt.Sprintf("%s:%d", myOpts.redisCfg.Ip, myOpts.redisCfg.Port), myOpts.redisCfg.DB)
	natsmq.Init(myOpts.natsCfg.Ip, myOpts.natsCfg.Port)
	senderId = serverId
	var rfaddr1 ServerRfAddr
	register = registry.Registery(&rfaddr1)

	multiConfig, err := configCs.ReaderConfig(myOpts.consulCfg.Ip, myOpts.consulCfg.Port, []string{"serverRegistry", _const.ST_MULTI})
	if err != nil {
		log.Println(err)
	}
	serverConfigs = make(map[string][]configCs.ServerRegistry)
	mJ := make([]configCs.ServerRegistry, 1)
	json.Unmarshal(multiConfig, &mJ)
	serverConfigs[_const.ST_MULTI] = mJ

	//消息中间件订阅
	natsmq.AsyncNats(serverId, serverId, handleNatsMsg)
}

func SetRedisCfg(redisCfg configCs.RedisConfig) opt {
	return func(o *opts) {
		fmt.Println(util.RunFuncName(), redisCfg)
		o.redisCfg = redisCfg
	}
}

func SetNatsCfg(natsCfg configCs.NatsConfig) opt {
	return func(o *opts) {
		fmt.Println(util.RunFuncName(), natsCfg)
		o.natsCfg = natsCfg
	}
}

func SetConsulCfg(consulCfg configCs.ConsulConfig) opt {
	return func(o *opts) {
		fmt.Println(util.RunFuncName(), consulCfg)
		o.consulCfg = consulCfg
	}
}

func SetRabbitmqCfg(rabbitmqCfg configCs.RabbitmqConfig) opt {
	return func(o *opts) {
		fmt.Println(util.RunFuncName(), rabbitmqCfg)
		o.rabbitmqCfg = rabbitmqCfg
	}
}
