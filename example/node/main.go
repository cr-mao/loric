/**
User: cr-mao
Date: 2023/8/23 11:35
Email: crmao@qq.com
Desc: main.go
*/
package main

import (
	"math/rand"
	"time"

	"github.com/cr-mao/loric"
	"github.com/cr-mao/loric/cluster/node"
	"github.com/cr-mao/loric/component"
	"github.com/cr-mao/loric/conf"
	"github.com/cr-mao/loric/example/node/router"
	"github.com/cr-mao/loric/locate/redis"
	"github.com/cr-mao/loric/registry/etcd"
	"github.com/cr-mao/loric/transport/grpc"
)

func main() {
	conf.InitConfig("local")
	//随机数种子
	rand.Seed(time.Now().UnixNano())
	// 配置初始化，依赖命令行 --env 参数
	//全局设置时区
	var cstZone, _ = time.LoadLocation(conf.GetString("app.timezone"))
	time.Local = cstZone
	contanier := loric.NewContainer()
	location := redis.NewLocator()
	serverOpts := make([]grpc.ServerOption, 0, 1)
	serverOpts = append(serverOpts,
		grpc.WithUnaryInterceptor(
			grpc.UnaryCrashInterceptor,
		),
	)
	rpcServer := node.NewTransport(node.WithServerOptions(
		serverOpts...,
	))
	nodeInstance := node.NewNode(
		node.WithLocator(location),
		node.WithTransporter(rpcServer),
		node.WithRegistry(etcd.NewRegistry()),
	)

	// 注册路由
	router.NewRouter(nodeInstance.Proxy()).Init()

	// 添加网关组件, pprof分析
	contanier.Add(nodeInstance, component.NewPProf())
	// 启动容器
	contanier.Serve()
}
