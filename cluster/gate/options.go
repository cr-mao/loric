package gate

import (
	"context"
	"time"

	"github.com/cr-mao/loric/conf"
	"github.com/cr-mao/loric/locate"
	"github.com/cr-mao/loric/network"
	"github.com/cr-mao/loric/registry"
	"github.com/cr-mao/loric/sugar"
	"github.com/cr-mao/loric/transport/grpc"
)

const (
	defaultName    = "gate"          // 默认名称
	defaultTimeout = 3 * time.Second // 默认超时时间
)

const (
	defaultIDKey      = "cluster.gate.id"
	defaultNameKey    = "cluster.gate.name"
	defaultTimeoutKey = "cluster.gate.timeout"
)

type Option func(o *options)

type options struct {
	id                string              // 实例ID
	name              string              // 实例名称
	ctx               context.Context     // 上下文
	timeout           time.Duration       // RPC调用超时时间
	server            network.Server      // 网关服务器
	locator           locate.Locator      // 用户定位器
	registry          registry.Registry   // 服务注册器
	grpcServerOptions []grpc.ServerOption // grpcServer 选项
	rpcAddr           string
}

func defaultOptions() *options {
	opts := &options{
		ctx:     context.Background(),
		name:    defaultName,
		timeout: defaultTimeout,
		rpcAddr: ":0",
	}

	if id := conf.GetString(defaultIDKey); id != "" {
		opts.id = id
	} else {
		if uuid, err := sugar.UUID(); err != nil {
			opts.id = uuid
		}
	}

	if name := conf.GetString(defaultNameKey); name != "" {
		opts.name = name
	}

	if timeout := conf.GetInt64(defaultTimeoutKey); timeout > 0 {
		opts.timeout = time.Duration(timeout) * time.Second
	}

	return opts
}

// WithID 设置实例ID
func WithID(id string) Option {
	return func(o *options) { o.id = id }
}

// WithName 设置实例名称
func WithName(name string) Option {
	return func(o *options) { o.name = name }
}

// WithContext 设置上下文
func WithContext(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

// WithServer 设置服务器
func WithServer(server network.Server) Option {
	return func(o *options) { o.server = server }
}

// WithTimeout 设置RPC调用超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) { o.timeout = timeout }
}

// WithLocator 设置用户定位器
func WithLocator(locator locate.Locator) Option {
	return func(o *options) { o.locator = locator }
}

// WithRegistry 设置服务注册器
func WithRegistry(r registry.Registry) Option {
	return func(o *options) { o.registry = r }
}

//  grpc server options
func WithGrpcServerOptions(grpcOptions ...grpc.ServerOption) Option {
	return func(o *options) { o.grpcServerOptions = grpcOptions }
}

func WithRpcAddr(rpcAddr string) Option {
	return func(o *options) { o.rpcAddr = rpcAddr }
}
