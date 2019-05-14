#!/bin/bash
## 这个脚本使用Docker在不同的网络命名空间产生多个client实例.
## 这样才能避免source port的限制，在一台机器上才能创建百万的连接.
##
## 用法: ./connect <connections> <number of clients> <server ip>
## Server IP 通常是 Docker gateway IP address, 缺省是 172.17.0.1

#CONNECTIONS=$1
#REPLICAS=$2
#IP=$3
CONNECTIONS=20000
REPLICAS=5
IP="172.17.0.1"
#go build --tags "static netgo" -o client client.go

echo $CONNECTIONS
echo $IP
echo $(pwd)

for((c=0; c < $REPLICAS; c++));
do
	docker container kill 1mclient_$c
	docker container rm 1mclient_$c
    #docker run -i -v $(pwd)/client:/client --name 1mclient_$c -d alpine /client -conn=$CONNECTIONS -ip=$IP
    #docker run -i -v $(pwd)/client:/client --name 1mclient_$c -d ubuntu /client -conn=$CONNECTIONS -ip=$IP
done
