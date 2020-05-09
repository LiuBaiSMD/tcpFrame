# tcp practice
## 1.发送普通字符串二进制数据进行输出

## 2.tcp认证
### ①建立tcp连接
### ②启动login登录认证
### ③认证失败关闭连接
### ④认证成功加入管理队列

## 3.管理队列（先整理一下代码） 
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
{
	Action string
	Name string
	PWD string
	UserId int
}

转换成

{
    LoadCode:1
    BytesData:[]byte
}
```
### ②通过标志位控制解析包