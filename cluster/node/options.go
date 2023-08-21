package node

import (
	"context"
	"time"

	"github.com/cr-mao/loric/conf"
	"github.com/cr-mao/loric/encoding"
	"github.com/cr-mao/loric/locate"
	"github.com/cr-mao/loric/registry"
	"github.com/cr-mao/loric/sugar"
)

type Option func(o *options)

type options struct {
	id          string            // 实例ID
	name        string            // 实例名称
	ctx         context.Context   // 上下文
	codec       encoding.Codec    // 编解码器
	timeout     time.Duration     // RPC调用超时时间
	locator     locate.Locator    // 用户定位器
	registry    registry.Registry // 服务注册器
	transporter *Transport        // 消息传输器
}

func defaultOptions() *options {
	opts := &options{
		ctx:     context.Background(),
		name:    defaultName,
		codec:   encoding.Invoke(defaultCodec),
		timeout: defaultTimeout,
	}

	if id := conf.GetString(defaultIDKey, ""); id != "" {
		opts.id = id
	} else if id, err := sugar.UUID(); err == nil {
		opts.id = id
	}

	if name := conf.GetString(defaultNameKey, ""); name != "" {
		opts.name = name
	}

	if codec := conf.GetString(defaultCodecKey); codec != "" {
		opts.codec = encoding.Invoke(codec)
	}

	if timeout := conf.GetInt64(defaultTimeoutKey, 0); timeout > 0 {
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

// WithCodec 设置编解码器
func WithCodec(codec encoding.Codec) Option {
	return func(o *options) { o.codec = codec }
}

// WithContext 设置上下文
func WithContext(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

// WithTimeout 设置RPC调用超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) { o.timeout = timeout }
}

// WithLocator 设置定位器
func WithLocator(locator locate.Locator) Option {
	return func(o *options) { o.locator = locator }
}

// WithRegistry 设置服务注册器
func WithRegistry(r registry.Registry) Option {
	return func(o *options) { o.registry = r }
}

// WithTransporter 设置消息传输器
func WithTransporter(transporter *Transport) Option {
	return func(o *options) { o.transporter = transporter }
}
