/**
User: cr-mao
Date: 2023/8/19 09:49
Email: crmao@qq.com
Desc: transport.go
*/
package gate

import (
	"google.golang.org/grpc"

	"github.com/cr-mao/loric/conf"
	"github.com/cr-mao/loric/internal/endpoint"
	mygrpc "github.com/cr-mao/loric/transport/grpc"
	"github.com/cr-mao/loric/transport/grpc/pb"
)

type ServerOptions struct {
	Addr       string
	ServerOpts []mygrpc.ServerOption
}

type transOptions struct {
	serverOptions ServerOptions        // gate grpc server
	clientOptions mygrpc.ClientOptions // 请求 node 的 grpc client 选项
}

type TransOptionFunc func(o *transOptions)

type Transport struct {
	opts *transOptions
	//once          sync.Once //  ClientBuilder 只用一次
	clientBuilder *mygrpc.ClientBuilder
}

func defaultTransOptions() *transOptions {
	opts := &transOptions{}
	opts.serverOptions.Addr = conf.Get(defaultGrpcServerAddrKey, defaultGrpcServerAddr)
	opts.clientOptions.PoolSize = conf.GetInt(defaultGrpcClientPoolSizeKey, defaultGrpcClientPoolSize)
	return opts
}

// WithServerListenAddr 设置服务器监听地址
func WithServerListenAddr(addr string) TransOptionFunc {
	return func(o *transOptions) { o.serverOptions.Addr = addr }
}

// WithServerOptions 设置服务器选项
func WithServerOptions(opts ...mygrpc.ServerOption) TransOptionFunc {
	return func(o *transOptions) { o.serverOptions.ServerOpts = opts }
}

// WithClientPoolSize 设置客户端连接池大小
func WithClientPoolSize(size int) TransOptionFunc {
	return func(o *transOptions) { o.clientOptions.PoolSize = size }
}

// WithClientDialOptions 设置客户端拨号选项
func WithClientDialOptions(opts ...grpc.DialOption) TransOptionFunc {
	return func(o *transOptions) { o.clientOptions.DialOpts = opts }
}

func NewTransport(opts ...TransOptionFunc) *Transport {
	trOpts := defaultTransOptions()
	for _, opt := range opts {
		opt(trOpts)
	}
	return &Transport{
		opts:          trOpts,
		clientBuilder: mygrpc.NewClientBuilder(&trOpts.clientOptions),
	}
}

func (t *Transport) NewGateServer(provider *provider) (*mygrpc.Server, error) {
	s, err := mygrpc.NewServer(t.opts.serverOptions.Addr, t.opts.serverOptions.ServerOpts...)
	if err != nil {
		return nil, err
	}
	pb.RegisterGateServer(s, &rpcService{
		provider: provider,
	})
	return s, nil
}

func (t *Transport) NewNodeClient(ep *endpoint.Endpoint) (*NodeGrpcClient, error) {
	cc, err := t.clientBuilder.GetConn(ep.Target())
	if err != nil {
		return nil, err
	}
	return NewNodeClient(cc), nil
}
