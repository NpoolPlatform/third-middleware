# Npool third middleware

[![Test](https://github.com/NpoolPlatform/third-middleware/actions/workflows/main.yml/badge.svg?branch=master)](https://github.com/NpoolPlatform/third-middleware/actions/workflows/main.yml)

[目录](#目录)
- [功能](#功能)
- [命令](#命令)
- [步骤](#步骤)
- [最佳实践](#最佳实践)
- [关于mysql](#关于mysql)

-----------

### 命令
* make init ```初始化仓库，创建go.mod```
* make verify ```验证开发环境与构建环境，检查code conduct```
* make verify-build ```编译目标```
* make test ```单元测试```
* make generate-docker-images ```生成docker镜像```
* make third-middleware ```单独编译服务```
* make third-middleware-image ```单独生成服务镜像```
* make deploy-to-k8s-cluster ```部署到k8s集群```

### 步骤
* 在github上将模板仓库https://github.com/NpoolPlatform/third-middleware.git import为https://github.com/NpoolPlatform/my-service-name.git
* git clone https://github.com/NpoolPlatform/my-service-name.git
* cd my-service-name
* mv cmd/third-middleware cmd/my-service
* 修改cmd/my-service/main.go中的serviceName为My Service
* mv cmd/my-service/ServiceTemplate.viper.yaml cmd/my-service/MyService.viper.yaml
* 将cmd/my-service/MyService.viper.yaml中的内容修改为当前服务对应内容
* 修改Dockerfile和k8s部署文档为当前服务对应内容
  * grep -rb "service template" ./*
  * grep -rb "ServiceTemplate" ./*
  * grep -rb "Service Template" ./*
  * grep -rb "service_template" ./*
  * grep -rb "third-middleware" ./*
  * grep -rb "servicetmpl" ./*
  * 修改cmd/my-service/k8s中的三个yaml文件，包含端口，服务名字
  * grep -rb "serviceid"，并使用uuid生成新值替换

### 最佳实践
* 每个服务只提供单一可执行文件，有利于docker镜像打包与k8s部署管理
* 每个服务提供http调试接口，通过curl获取调试信息
* 集群内服务间direct call调用通过服务发现获取目标地址进行调用
* 集群内服务间event call调用通过rabbitmq解耦

### 关于mysql
* 创建app后，从app.Mysql()获取本地mysql client
* [文档参考](https://entgo.io/docs/sql-integration)