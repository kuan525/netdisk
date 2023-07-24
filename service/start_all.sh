#!/bin/bash

# 检查service进程
check_process() {
    sleep 1
    res=`ps aux | grep -v grep | grep "service/bin" | grep $1`
    if [[ $res != '' ]]; then
        echo -e "\033[32m 已启动 \033[0m" "$1"
        return 1
    else
        echo -e "\033[31m 启动失败 \033[0m" "$1"
        return 0
    fi
}

# 编译service可执行文件, 在netdisk目录下
build_service() {
    cd service/$1
    go build
    mv $1 ../bin
    cd ../..
    resbin=`ls service/bin/ | grep $1`
    echo -e "\033[32m 编译完成: \033[0m service/bin/$resbin"

}

# 启动service
run_service() {
    ./service/bin/$1 --registry=consul >> $logpath/$1.log 2>&1 &
    sleep 1
    check_process $1
}

# 创建运行日志目录
logpath=/Users/kuan525/TODO/log
mkdir -p $logpath

# 切换到工程根目录
#cd $GOPATH/netdisk
#cd /data/go/work/src/netdisk

# 微服务可以用supervisor做进程管理工具；
# 或者也可以通过docker/k8s进行部署

services="
dbproxy
upload
download
transfer
account
apigw
"

# 启动rabbitmq
#rabbitmq-server
# 启动consul
#consul agent -dev

# 检测端口是否被使用
check_port() {
    port=$1
    if lsof -i :$port >/dev/null; then
        echo "Port $port is in use. Closing the port..."
        sudo kill $(sudo lsof -t -i :$port)
        echo "Port $port has been closed."
    else
        echo "Port $port is not in use."
    fi
}

# 检测并关闭端口
check_port 8080
check_port 28080
check_port 38080

# 执行编译service
mkdir -p service/bin/ && rm -f service/bin/*
for service in $services
do
    build_service $service
done

# 执行启动service
for service in $services
do
    run_service $service
done

echo '微服务启动完毕.'

