# tcp frame
# 项目功能
```
进行tcp连接，并进行消息监听、转发，存储至redis等发布订阅的结构中;
通过发布订阅与其他的服务（例如微服务、或者自定义进程服务）进行交互；
暂定使用rabbitmq；

```

# 项目结构
```
├── client      #客户端连接代码
├── config      #配置中心
├── conns       #连接管理
├── const       #常量存储
├── dao         #数据库中心
├── datas       #数据结构管理
├── handle      #处理方法
├── msg         #tcp消息监听、解析、分发
├── registry    #处理方法注册中心
├── server      #服务端启动代码
└── util        #自定义工具
```

## 1.发送普通字符串二进制数据进行输出

## 2.tcp认证
### ①建立tcp连接
### ②启动login登录认证
### ③认证失败关闭连接
### ④认证成功加入管理队列

## 3.管理队列
### ①加入心跳机制
```
成功登陆后，开启心跳检测
```
### ②掉线断开连接
```
掉线关闭程序，使用管道同学关闭信号
```
### ③断线重连机制
```
1.模拟建立通信
2.在通信过程中断开连接
3.期间不断建立通信
4.断开的client带着userId重连
5.恢复之前的通信
```
### ④断线后从连接队列中删除

## 4.心跳管理 
### ①管理连接的时间戳，建立连接的时间
### ②超时将删除连接
```
超过三次心跳没有发送心跳包的连接将会被关闭
```

## 5.改进包协议
### ①在打包协议时将标志位进行打包
```
data = {
	Action string
	Name string
	PWD string
	UserId int
}

转换成二进制数据

bData = []byte(data)

{
    LoadCode:1
    BytesData:bData
}
```
### ②通过标志位控制解析包

### ③将方法注册到msg-registry中

### ④自动根据action，找到对应的msg-registry，进行处理

## 6.改进数据包传输协议
### ①增加组装传输数据的接口
```
总共分为两层
1.(第一层解析)数据包长度dataLenth（32位 []byte）+ 编码类型codeType（8位 []byte）+ 数据data（[]byte）
dataLenth:存储data长度
codeType:基础解析格式，标识解析data的方式，json、proto等通用的格式
data:数据内容

2.(第二层解析)解析data模块，将data分解成各个类型json、proto等的BaseData后，其中的Action数据为指导业务层自行解析的模块，比如
例① json中的BaseData结构:
type BaseData struct{
    Action string,
    UserId int,
    BData []byte,
}

json中的HeartBeat结构:
type HeartBeat struct{
    Action string,
    UserId int,
    TimeStamp int,
    OtherMsg string,
}

例如在上述拆包过程中codeType=1代表json格式数据，将data解析为json的BaseData格式:得到以下数据
data = {
        Action:"Heartbeat",
        UserId:10001,
        BData:[12, 23, 45, 234, 54, 65],
        }
(第三层解析)然后业务层通过Action将指导BData解析为已经定义好的json结构 HeartBeat
BData = {
    Action: HeartBeat,
    UserId: 10001,
    TimeStamp: 123456789,
    OtherMsg: "hello world!",
}

例②
proto中的BaseData结构:
message BaseData {
    string Action = 1;
    int64 UserId = 2;
    bytes BData = 3;
}

proto中的HeartBeat结构:
message HeartBeat {
    string Action = 1;
    int64 UserId = 2;
    int64 TimeStamp = 3;
    string OtherMsg = 4;
}

例如在上述拆包过程中codeType=2代表proto格式数据，将data解析为proto的BaseData格式:得到以下数据
data = {
        Action:"Heartbeat",
        UserId:10001,
        BData:[12, 23, 45, 234, 54, 65],
        }
(第三层解析)然后业务层通过Action将BData解析为已经定义好的proto结构 HeartBeat，
BData = {
    Action: HeartBeat,
    UserId: 10001,
    TimeStamp: 123456789,
    OtherMsg: "hello world!",
}

```
