/*
@Author: liubai
@Date: 2020/5/16 6:37 下午
@Desc: use for what
*/

package msgMQ


var RabbitMq RabbitMQAMQP
func init(){
	RabbitMq = RabbitMQAMQP{
		rbtname: "guest",
		passwd: "guest",
		ipAddr: "localhost",
		port: 5672,
		exchangeMap:make(map[string]ExchangeAMQP),
	}
	RabbitMq.connect()
}
