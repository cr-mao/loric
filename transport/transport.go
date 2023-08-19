/**
User: cr-mao
Date: 2023/8/17 14:03
Email: crmao@qq.com
Desc: 传输器接口
*/
package transport

import (
	"github.com/cr-mao/loric/internal/endpoint"
)

type Server interface {
	// Start 启动服务器
	Start() error
	// Stop 停止服务器
	Stop() error
	//服务地址
	Endpoint() *endpoint.Endpoint
	// Addr 监听地址
	Addr() string
	// Scheme 协议
	Scheme() string
}
