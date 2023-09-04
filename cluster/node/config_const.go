/**
User: cr-mao
Date: 2023/8/19 20:16
Email: crmao@qq.com
Desc: config_const.go
*/
package node

import "time"

const (
	defaultIDKey                 = "node.id"
	defaultNameKey               = "node.name"
	defaultCodecKey              = "node.codec"
	defaultTimeoutKey            = "node.timeout"
	defaultWeightKey             = "node.weight"
	defaultGrpcServerAddrKey     = "node.grpc.server.addr"
	defaultGrpcClientPoolSizeKey = "node.grpc.client.poolSize"
)

const (
	defaultName               = "node_01"       //  应用名称
	defaultTimeout            = 3 * time.Second // 默认超时时间
	defaultWeight             = 10              //默认服务权重
	defaultCodec              = "proto"         // 默认编解码器名称
	defaultGrpcServerAddr     = ":0"            //  默认服务器地址
	defaultGrpcClientPoolSize = 10              // 默认请求node grpc的客户端连接池大小
)
