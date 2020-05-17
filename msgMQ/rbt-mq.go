/*
@Author: liubai
@Date: 2020/5/16 6:37 下午
@Desc: use for what
*/

package msgMQ

var (
	RabbitMQMap map[string]*RabbitMQAMQP //使用serviceName， RabbitMq作为数据存储
)

func init(){
	RabbitMQMap = make(map[string]*RabbitMQAMQP)
	err:=NewRabbitMQAMQP("server1", "guest", "guest", "localhost", 5672)
	if err!=nil{
		panic(err)
	}
}

func NewRabbitMQAMQP(rbtServiceName, rbtname, passwd, ipAddr string, port int)error{
	rabbitAMQP := RabbitMQAMQP{
		rbtname: rbtname,
		passwd: passwd,
		ipAddr: ipAddr,
		port: port,
		exchangeMap:make(map[string]ExchangeAMQP),
	}
	err := rabbitAMQP.connect()
	RabbitMQMap[rbtServiceName] = &rabbitAMQP
	return err
}