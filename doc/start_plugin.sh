# 启动consul服务
consul agent -dev &
# 启动rabbitmq服务
rabbitmq-server &
# 启动nats服务
nats-server &
# 启动redis服务
redis-server &