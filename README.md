# tcp frame
# 项目功能
```
进行tcp连接，并进行消息监听、转发，存储至redis等发布订阅的结构中;
通过发布订阅与其他的服务（例如微服务、或者自定义进程服务）进行交互；
暂定使用rabbitmq；

```

# 项目结构
```
.
├── client          #客户端连接代码
├── config          #配置中心
├── conns           #连接管理
├── const           #常量存储
├── dao             #数据库中心
├── datas           #数据结构管理
├── doc             #文档说明
├── handle          #处理方法
├── msg             #tcp消息监听、解析、分发
├── msgMQ           #消息中间件
├── registry        #处理方法注册（类似于装饰器）
├── server          #客户端使用代码
├── server-registry #服务注册
└── util            #使用的工具
```

# 系统结构图
![avatar](https://github.com/LiuBaiSMD/tcpFrame/blob/master/doc/TcpFrame.jpg?raw=true)

# 系统功能
## 1.建立用户与服务端tcp连接，自行解析包
```
①使用protobuf，对用户发送的tcp连接进行序列化反序列化处理
②包的解析：
 tcp中传输的二进制数据可以概括为
 头部长度（四个字节）+ 消息长度（四个字节）+ 头部二进制数据(proto格式RequestHeader) + 消息二进制数据(proto格式TokenTcpRequest)
```
## 2.进行用户验证，管理用户tcp连接
```
①用户连接后，通过http请求获得token通行证
②用户使用token通行证，在连接tcp后，发送第一次请求时，请求token方式连接
③tcpConn服通过验证后对连接进行统一管理，并设置生命周期，服务端需要自行发送心跳包heartBeat
④通过心跳包进行用户管理，超时连接将关闭
```

## 3.微服务集群
```
①用户通过tcp连接后，只需要发送服务的名字serverType以及指令名称cmdType,便可请求对应的服务；
②服务的集群可以自由创建删除，并且对用户的连接无任何影响
```
```
微服务搭建：
微服务参考 clinet/server/token-server
1.连接对应的consul，进行consul服务连接
server-registry.ConsulConnect("localhost:8500")
2.进行服务注册，注册服务的ip、端口、服务名字、tags
serverId, _ = server-registry.RegisterServer(
	"127.0.0.1",
	0,
	serverName,
	[]string{})
3.订阅消息频道（消息格式MsgBody， 通过其中的CmdType供业务自行解析）
订阅自己serverName的频道：接收公共服务的消息
go natsmq.AsyncNats(serverName, serverName, handleNatsMsg)
订阅自己serverId的频道：接收定向消息
go natsmq.AsyncNats(serverId, serverId, handleNatsMsg)

handleNatsMsg为微服务中需要接收到消息后处理的方法
```

## 4.配置管理
```
连接consul后
使用config/consul中的方法，上传固定文件到指定路径
使用config/consul中的方法，获取指定路径的配置数据
```

## 5.消息中间件
```
消息中间件有rabbitmq、nats两种，由于业务暂时比较简单，目前只需使用nats即可
    go natsmq.AsyncNats(serverName, workGroup, handleNatsMsg)
    func handleNatsMsg(msg *nats.Msg) {}
监听serverName频道，加入workGroup组，并将消息自动加入handlleNatsMsg处理

```

## 6.数据传输格式
```
使用proto数据进行传输，各个服务需要在使用时自行定义服务名字serverName，以及CmdType对应的proto结构体
并使用 proto --proto_path=. --go_out=. you_proto_file_name.proto，生成对应.go文件
```

## 7.registry模块，方法注册使用
```
registry使用反射的方式将对应的方法注册成map[funcName]func，使用参考registry test模块
```

# 使用教程

## 安装环境
```
安装nats
安装consul
安装rabbitmq
安装redis
安装mongo

在主目录下 ./tcpFrame中执行
go mod tidy
```
## 启动插件服务
```
请确保使用默认ip、端口 todo 使用配置构造
sh ./tcpFrame/doc/start_plugin.sh
```

## 启动tcp链接服，tcpConn主服务，server/server.go
```
使用go mod tidy下载依赖包
1.启动rabbitmq 本地启动rabbitmq 使用默认端口
2.启动consul 本地启动consul ： consul agent -dev
3.启动go run ./tcpFrame/server/server.go
```

## 启动http-token token管理http服务
```
go run ./tcpFrame/client/http-token/http-token.go
```

## 启动client 模拟用户请求
```
1.启动client中的client，模拟客户端请求
go run clinet
2.启动client中的token-server,模拟服务集群中的token生成服务
go run ./tcpFrame/client/token-server.go
```