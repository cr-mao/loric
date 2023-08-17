package packet

import (
	"encoding/binary"

	"github.com/cr-mao/loric/conf"
)

const (
	defaultMessageBytes = 5000
)

const (
	defaultMessageBytesKey = "packet.bufferBytes"
)

// -------------------------
// | route | seq | message |
// -------------------------
type options struct {
	// 字节序
	// binary.BigEndian 写死
	byteOrder binary.ByteOrder
	// 路由字节数
	// 写死2字节
	routeBytes int
	// 序列号字节数
	// 写死4字节
	seqBytes int
	// 消息字节数
	// 默认为5000字节
	bufferBytes int
}

type Option func(o *options)

func defaultOptions() *options {
	opts := &options{
		byteOrder:   binary.BigEndian,
		routeBytes:  2,
		seqBytes:    4,
		bufferBytes: conf.GetInt(defaultMessageBytesKey, defaultMessageBytes),
	}
	return opts
}

// WithByteOrder 设置字节序
//func WithByteOrder(byteOrder binary.ByteOrder) Option {
//	return func(o *options) { o.byteOrder = byteOrder }
//}
//// WithSeqBytes 设置序列号字节数
//func WithSeqBytes(seqBytes int) Option {
//	return func(o *options) { o.seqBytes = seqBytes }
//}

// WithMessageBytes 设置消息字节数
func WithMessageBytes(messageBytes int) Option {
	return func(o *options) { o.bufferBytes = messageBytes }
}
