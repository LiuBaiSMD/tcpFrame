/*
@Author: liubai
@Date: 2020/5/31 4:28 下午
@Desc: use for what
*/

package consul

type ServerRegistry struct {
	Name string `json:"name"`
	Tags []string `json:"tags"`
	RouteKey string `json:"routeKey"`
}

type ConsulConfig struct {
	Ip string `json:"Ip"`
	Port int `json:"Port"`
}

type RabbitmqConfig struct {
	ServerName string `json:"ServerName"`
	Ip string `json:"Ip"`
	Port int `json:"Port"`
	LoginName string `json:"LoginName"`
	Password string `json:"Password"`
}