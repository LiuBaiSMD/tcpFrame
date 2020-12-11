# 启动consul服务
ps -ef | grep consul | grep -v "grep" | awk '{print $2}'| xargs kill # 关闭rabbitmq服务
sudo rabbitmqctl stop
ps -ef | grep nats | grep -v "grep" | awk '{print $2}' | xargs kill
redis-cli shutdown