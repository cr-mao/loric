/**
User: cr-mao
Date: 2023/8/17 14:03
Email: crmao@qq.com
Desc: 传输器接口
*/
package transport

import (
	"context"
	"net/url"
)

// Kind defines the type of Transport
type Kind string

func (k Kind) String() string { return string(k) }

const (
	KindGRPC Kind = "grpc"
	KindHttp Kind = "http"
	KindTcp  Kind = "tcp"
)

type Server interface {
	// Start 启动服务器
	Start(ctx context.Context) error
	// Stop 停止服务器
	Stop(ctx context.Context) error
}

// Endpointer is registry endpoint.
type Endpointer interface {
	Endpoint() (*url.URL, error)
}
