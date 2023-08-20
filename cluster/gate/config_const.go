/**
User: cr-mao
Date: 2023/8/19 10:06
Email: crmao@qq.com
Desc: constant.go
*/
package gate

import "time"

const (
	defaultIDKey      = "gate.id"      // 实例id 配置key
	defaultNameKey    = "gate.name"    // 应用名称 配置key
	defaultTimeoutKey = "gate.timeout" // grpc请求超时时间

	defaultGrpcServerAddrKey     = "gate.grpc.server.addr"
	defaultGrpcClientPoolSizeKey = "gate.grpc.client.poolSize"
)

const (
	defaultName    = "gate_01"       //  应用名称
	defaultTimeout = 3 * time.Second // 默认超时时间

	defaultGrpcServerAddr     = ":0" //  默认服务器地址
	defaultGrpcClientPoolSize = 10   // 默认请求node grpc的客户端连接池大小
)
