# 计划进度表

## 5.13计划
```
1.拆分各个模块功能，进行模块独立化设计
2.模块功能实现
3.完成链接层的数据拉取、链接管理梳理
4.完成消息队列连接，发布订阅的使用
```
## 连接队列改进
```
1.tcp 长链接（开启keepalive 心跳 ） 建立长链接并维护长链接关系表
2.断开（正常断开和异常断开） 更新长链接维护表、更新登录队列表（有可能更新游戏人数）

3.登录 登录请求 -》 获取登录编号-〉入登录队列 （server协程处理 queue chan、位置编号管理 map）处理到哪一个就返回哪一个
（在返回之前，客户端可以实时来查询自己在队列中的位置）
4.提供查询队列位置接口
5.提供退出队列接口
6.网络断开处理
```

## 5.17计划
```
完成任务队列的监听
consumerabbit->queueName
监听前监测队列是否存在，定义消费者
如果已经存在消费者key则取消前面的，一个连接之间不能声明同一个consumer
```

## 5.19完成
```
①完成msg模块使用json数据交互
②增加datas/proto的数据格式
```

## 5.19计划
```
1.增加专门的解释器通过cmdNo解析[]bytes为具体的proto格式
2.开始测试使用rabbitmq、查看nats
```

## 5.21记录
[protobuf存储原理概要介绍](https://blog.csdn.net/weixin_34029949/article/details/91461766) 
```
1.如果有多个字节连续的小头字节序,翻转【字节】顺序
2.一个rabbitmq可以同时多个连接
3.每个连接开通一个channel()，会开辟一条通道，互不干扰consumer
4.同一个channel不能注册同一个消费者
```

## 5.21计划
```
封装rabbitmq的消息方法
```
## 5.22计划
```
完成consul配置系统以及服务注册的接入
进行简单地服务注册
```
## 5.23计划
```
接入consul配置系统
功能
1.通过指定文件路径、指定ip上传配置
2.文件监听，在文件变动时更新consul上的配置
3.获取配置，传入对应的url以及路径获取对应的字段配置
```
## 5.24记录
```
consul中一个服务只能注册一次，服务名字不能重复
1.服务类型名字
2.服务的Id
3.tag 服务版本
增加config、server_regitry模块
```

## 5.24计划
```
联合rabbitmq、server_registry模块使用
```

## 5.25记录
```
发现rabbitmq中检测队列存在以及交换机存在的方法不正常返回
待修复注册queue、exchange错误
```

## 5.27计划
```
规定服务建立后，监听的信息接收queuename频道以及发送的queuename
```

## 5.28计划
```
需要解决的问题：服务启动时注册，监听两个消息队列
消息队列：一种是只要服务名字一样，都可以机会均等的接受到任务，；
第二种：可能有的用户信息是确定的，指定要serverId,或者不同的服，服务启动时注册一个直连的队列

尝试修改队列属性，在没有消费者或者生产者是删除队列


链接服在监听到服务启动后：
链接服启动：注册两个rabbitmq
一个queueName=serverName+server，的一个rbt队列
一个queueName=serverName+serverId，的一个定向rbt队列

普通服务启动：
每个服务启动时，需要注册serverName、ip、port、tags到consul中
监听queueName=serverName+server，的一个rbt队列
监听queueName=serverName+serverId，的一个定向rbt队列

链接服监听到普通服务启动后：
如何监听：链接服监听一个consul上的配置，每一个服务启动时，将会将自己的server、serverId等信息存入到consul配置中；
链接服根据配置或者定时检测consul中注册的服务，根据服务启动的server信息发布两个rabbitmq队列：
一个queueName=serverName+server，的一个rbt队列
一个queueName=serverName+serverId，的一个定向rbt队列

```

## 5.29计划
```
服务通过消息发布到固定的服务,第一个使用serverType
消息解析，配置读取
消息注册实现单例进程以及多实例进程
```

## 6.1计划
```
1.增加单独的配置上传服务
2.增加服务内配置获取插件
3.完成客户端模拟的消息获取过程

使用nats作为消息订阅插件，启用回调机制，当消息队列收到消息时自动调用对应的方法
回调方法中处理消息的基础解析，并启动对应的处理方法
```

## 6.1记录
```
nats封装记录：
使用serverName+"req"作为queueName作为请求数据，
使用serverName+"rsq"作为queueName作为回复数据，
```

## 6.3.1计划
```
datas/MsgBody
//此结构为消息中间件中的数据传输格式，其中cmdType供服务中自行区分解析
//加上发送方的sender_id, 接受的userId
```

## 6.3.2计划
```
datas/MsgBody
//此结构为消息中间件中的数据传输格式，其中cmdType供服务中自行区分解析
//加上发送方的sender_id, 接受的userId
```

## 6.19计划
```
改进tcp客户端，防止出现short write
```

## 3.08计划
```
nats通信是否可以指定特定server以及通用server
```