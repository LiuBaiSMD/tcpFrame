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