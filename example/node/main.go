/**
User: cr-mao
Date: 2023/8/23 11:35
Email: crmao@qq.com
Desc: main.go
*/
package main

import (
	"flag"
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

	nodeId := flag.String("node_id", conf.Get("node.id"), "节点id")
	nodeName := flag.String("node_name", conf.Get("node.name"), "节点名")
	pprofAddr := flag.String("pprof_addr", conf.Get("app.pprof.addr"), "pprof地址")
	grpcAddr := flag.String("grpc_addr", conf.GetString("node.grpc.server.addr"), "grpc addr")

	flag.Parse()
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
	),
		node.WithServerListenAddr(*grpcAddr),
	)
	nodeInstance := node.NewNode(
		node.WithID(*nodeId),
		node.WithName(*nodeName),
		node.WithLocator(location),
		node.WithTransporter(rpcServer),
		node.WithRegistry(etcd.NewRegistry()),
	)

	// 注册路由
	router.NewRouter(nodeInstance.Proxy()).Init()

	// 添加网关组件, pprof分析
	contanier.Add(nodeInstance, component.NewPProf(*pprofAddr))
	// 启动容器
	contanier.Serve()
}
