/*
@Author: liubai
@Date: 2020/5/31 4:28 下午
@Desc: use for what
*/

package consul

import (
	"encoding/json"
	"errors"
)

type ServerRegistry struct {
	Name     string   `json:"name"`
	Tags     []string `json:"tags"`
	RouteKey string   `json:"routeKey"`
}

type ConsulConfig struct {
	Ip   string `json:"Ip"`
	Port int    `json:"Port"`
}

type NatsConfig struct {
	Ip   string `json:"Ip"`
	Port int    `json:"Port"`
}

type RabbitmqConfig struct {
	ServerName string `json:"ServerName"`
	Ip         string `json:"Ip"`
	Port       int    `json:"Port"`
	LoginName  string `json:"LoginName"`
	Password   string `json:"Password"`
}

type RedisConfig struct {
	Ip       string `json:"Ip"`
	Port     int    `json:"Port"`
	DB       int    `json:"DB"`
	Password string `json:"Password"`
}

type TcpConn struct {
	Ip   string `json:"Ip"`
	Port int    `json:"Port"`
}

func GetRedisCfg(consulIp string, consulPort int) (*RedisConfig, error) {
	cfg := &RedisConfig{}
	if rawCfg, err := ReaderConfig(consulIp, consulPort, []string{"plugin", "redis"}); err == nil {
		if len(rawCfg)<=4{
			return nil, errors.New("no redis config!")
		}
		err = json.Unmarshal(rawCfg, cfg)
	}

	return cfg, err
}

func GetConsulCfg(consulIp string, consulPort int) (*ConsulConfig, error) {
	cfg := &ConsulConfig{}
	if rawCfg, err := ReaderConfig(consulIp, consulPort, []string{"plugin", "consul"}); err == nil {
		if len(rawCfg)<=4{
			return nil, errors.New("no Consul config!")
		}
		err = json.Unmarshal(rawCfg, cfg)
	}
	return cfg, nil
}

func GetRabbitmqCfg(consulIp string, consulPort int) (*RabbitmqConfig, error) {
	cfg := &RabbitmqConfig{}
	if rawCfg, err := ReaderConfig(consulIp, consulPort, []string{"plugin", "rabbitmq"}); err == nil {
		if len(rawCfg)<=4{
			return nil, errors.New("no Rabbitmq config!")
		}
		err = json.Unmarshal(rawCfg, cfg)
	}
	return cfg, nil
}

func GetNatsCfg(consulIp string, consulPort int) (*NatsConfig, error) {
	cfg := &NatsConfig{}
	if rawCfg, err := ReaderConfig(consulIp, consulPort, []string{"plugin", "nats"}); err == nil {
		if len(rawCfg)<=4{
			return nil, errors.New("no Nats config!")
		}
		err = json.Unmarshal(rawCfg, cfg)
	}
	return cfg, nil
}
