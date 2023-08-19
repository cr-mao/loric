/**
User: cr-mao
Date: 2023/7/31
Time: 15:09
Desc: main.go
*/
package main

import (
	"math/rand"
	"time"

	"github.com/cr-mao/loric"
	"github.com/cr-mao/loric/cluster/gate"
	"github.com/cr-mao/loric/component"
	"github.com/cr-mao/loric/conf"
	"github.com/cr-mao/loric/locate/redis"
	"github.com/cr-mao/loric/network/tcp"
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
	rpcServer := gate.NewTransport(gate.WithServerOptions(
		serverOpts...,
	))
	gateServer := gate.NewGate(
		gate.WithServer(tcp.NewServer()),
		gate.WithLocator(location),
		gate.WithTransport(rpcServer),
		gate.WithRegistry(etcd.NewRegistry()),
	)
	// 添加网关组件, pprof分析
	contanier.Add(gateServer, component.NewPProf())
	// 启动容器
	contanier.Serve()
}
