package tcp

import (
	"time"

	"github.com/cr-mao/loric/conf"
)

const (
	defaultServerAddr               = ":3553"
	defaultServerMaxConnNum         = 5000
	defaultServerHeartbeatInterval  = 10
	defaultHandlerMsgAsync          = false
	defaultServerHeatSend           = false
	defaultServerOneSecondMaxMsgNum = 10
	defaultMaxMsgSize               = 1024
)

const (
	defaultServerAddrKey               = "network.tcp.server.addr"
	defaultServerMaxConnNumKey         = "network.tcp.server.maxConnNum"
	defaultServerHeartbeatIntervalKey  = "network.tcp.server.heartbeatInterval"
	defaultHandlerMsgAsyncKey          = "network.tcp.server.handlerMsgAsync"
	defaultServerHeatSendKey           = "network.tcp.server.serverHeatSend"
	defaultServerOneSecondMaxMsgNumKey = "network.tcp.server.oneSecondMaxMsgNum"
	defaultMaxMsgSizeKey               = "network.tcp.server.maxMsgSize"
)

type ServerOption func(o *serverOptions)

type serverOptions struct {
	addr               string        // 监听地址，默认0.0.0.0:3553
	maxConnNum         int           // 最大连接数，默认5000
	heartbeatInterval  time.Duration // 心跳检测间隔时间，默认10s
	handleMsgAsync     bool
	serverHeatSend     bool
	oneSecondMaxMsgNum int
	maxMsgSize         uint32
}

func defaultServerOptions() *serverOptions {
	return &serverOptions{
		addr:               conf.GetString(defaultServerAddrKey, defaultServerAddr),
		maxConnNum:         conf.GetInt(defaultServerMaxConnNumKey, defaultServerMaxConnNum),
		heartbeatInterval:  time.Duration(conf.GetInt(defaultServerHeartbeatIntervalKey, defaultServerHeartbeatInterval)) * time.Second,
		handleMsgAsync:     conf.GetBool(defaultHandlerMsgAsyncKey, defaultHandlerMsgAsync),
		serverHeatSend:     conf.GetBool(defaultServerHeatSendKey, defaultServerHeatSend),
		oneSecondMaxMsgNum: conf.GetInt(defaultServerOneSecondMaxMsgNumKey, defaultServerOneSecondMaxMsgNum),
		maxMsgSize:         uint32(conf.GetInt32(defaultMaxMsgSizeKey, defaultMaxMsgSize)),
	}
}

// WithServerListenAddr 设置监听地址
func WithServerListenAddr(addr string) ServerOption {
	return func(o *serverOptions) { o.addr = addr }
}

// WithServerMaxConnNum 设置连接的最大连接数
func WithServerMaxConnNum(maxConnNum int) ServerOption {
	return func(o *serverOptions) { o.maxConnNum = maxConnNum }
}

// WithServerHeartbeatInterval 设置心跳检测间隔时间
func WithServerHeartbeatInterval(heartbeatInterval time.Duration) ServerOption {
	return func(o *serverOptions) { o.heartbeatInterval = heartbeatInterval }
}

// 接收消息 单连接 同步还是异步处理， 默认同步
func WithServerHandlerMsgAsync(handlerMsgAsync bool) ServerOption {
	return func(o *serverOptions) { o.handleMsgAsync = handlerMsgAsync }
}

// 是否心跳给客户端
func WithServerHeatSend(serverHeatSend bool) ServerOption {
	return func(o *serverOptions) { o.serverHeatSend = serverHeatSend }
}

// 1秒最大能接收包单次数 默认 10次
func WithServerOneSecondMaxMsgNum(oneSecondMaxMsgNum int) ServerOption {
	return func(o *serverOptions) { o.oneSecondMaxMsgNum = oneSecondMaxMsgNum }
}

// 1秒最大能接收的msg包的大小 默认1mb
func WithServerMaxMsgSize(maxMsgSize uint32) ServerOption {
	return func(o *serverOptions) { o.maxMsgSize = maxMsgSize }
}
