package gate

import (
	"context"
	"time"

	"github.com/cr-mao/loric/conf"
	"github.com/cr-mao/loric/locate"
	"github.com/cr-mao/loric/network"
	"github.com/cr-mao/loric/registry"
	"github.com/cr-mao/loric/sugar"
)

type Option func(o *options)

type options struct {
	id          string            // 实例ID
	name        string            // 实例名称
	weight      int               // 权重
	authTimeOut time.Duration     // 认证最长时间，默认5秒，连上来的第一个包必须是auth包
	ctx         context.Context   // 上下文
	server      network.Server    // 网关服务器
	locator     locate.Locator    // 用户定位器
	registry    registry.Registry // 服务注册器
	timeout     time.Duration     // grpc调用超时时间
	transporter *Transport        // 消息传输器
}

func defaultOptions() *options {
	opts := &options{
		ctx:         context.Background(),
		name:        defaultName,
		timeout:     defaultTimeout,
		authTimeOut: defaultAuthTimeout,
		weight:      10,
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

	if authTimeout := conf.GetInt64(defaultAuthTimeoutKey); authTimeout > 0 {
		opts.timeout = time.Duration(authTimeout) * time.Second
	}

	if weight := conf.GetInt(defaultWeightKey); weight > 0 {
		opts.weight = weight
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

// 消息传输， gate server,node client
func WithTransport(t *Transport) Option {
	return func(o *options) { o.transporter = t }
}

func WithAuthTimeOut(t int64) Option {
	return func(o *options) { o.authTimeOut = time.Duration(t) * time.Second }
}

func WithWeight(weight int) Option {
	return func(o *options) { o.weight = weight }
}
