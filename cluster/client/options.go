package client

import (
	"context"

	"github.com/cr-mao/loric/conf"
	"github.com/cr-mao/loric/encoding"
	"github.com/cr-mao/loric/network"
	"github.com/cr-mao/loric/sugar"
)

const (
	defaultName  = "client" // 默认客户端名称
	defaultCodec = "proto"  // 默认编解码器名称
)

const (
	defaultIDKey    = "client.id"
	defaultNameKey  = "client.name"
	defaultCodecKey = "client.codec"
)

type Option func(o *options)

type options struct {
	id     string          // 实例ID
	name   string          // 实例名称
	ctx    context.Context // 上下文
	codec  encoding.Codec  // 编解码器
	client network.Client  // 网络客户端
}

func defaultOptions() *options {
	opts := &options{
		ctx:   context.Background(),
		name:  defaultName,
		codec: encoding.Invoke(defaultCodec),
	}

	if id := conf.Get(defaultIDKey); id != "" {
		opts.id = id
	} else {
		if uuidN, err := sugar.UUID(); err == nil {
			opts.id = uuidN
		}
	}

	if name := conf.Get(defaultNameKey); name != "" {
		opts.name = name
	}

	if codec := conf.Get(defaultCodecKey); codec != "" {
		opts.codec = encoding.Invoke(codec)
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

// WithClient 设置客户端
func WithClient(client network.Client) Option {
	return func(o *options) { o.client = client }
}

// WithContext 设置上下文
func WithContext(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}
