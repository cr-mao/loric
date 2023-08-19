package cluster

import (
	"context"
	"github.com/cr-mao/loric/conf"
	"github.com/cr-mao/loric/locate"
	"github.com/cr-mao/loric/sugar"
	"time"

	"github.com/cr-mao/loric/component"
	"github.com/cr-mao/loric/log"
	"github.com/cr-mao/loric/network"
	"github.com/cr-mao/loric/registry"
	"github.com/cr-mao/loric/session"
	"github.com/cr-mao/loric/transport/grpc"
)

type Gate struct {
	component.Base
	opts      *options
	ctx       context.Context
	cancel    context.CancelFunc
	proxy     *proxy
	instance  *registry.ServiceInstance
	rpcServer *grpc.Server
	session   *session.Session
}

func NewGate(opts ...Option) *Gate {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	g := &Gate{}
	g.opts = o
	g.proxy = newProxy(g)
	g.session = session.NewSession()
	g.rpcServer = newGrpcServer(&provider{g})
	g.ctx, g.cancel = context.WithCancel(o.ctx)
	return g
}

// Name 组件名称
func (g *Gate) Name() string {
	return g.opts.name
}

// Init 初始化
func (g *Gate) Init() {
	if g.opts.id == "" {
		log.Fatal("instance id can not be empty")
	}
	if g.opts.server == nil {
		log.Fatal("server component is not injected")
	}
	if g.opts.locator == nil {
		log.Fatal("locator component is not injected")
	}
	if g.opts.registry == nil {
		log.Fatal("registry component is not injected")
	}
	if len(g.opts.grpcServerOptions) == 0 {
		log.Fatal("grpcServerOptions is not injected")
	}
}

// Start 启动组件
func (g *Gate) Start() {
	g.startNetworkServer()

	g.startRPCServer()

	g.registerServiceInstance()

	// 启动监听定位器和node服务发现n
	g.proxy.watch(g.ctx)

	g.debugPrint()
}

// Destroy 销毁组件
func (g *Gate) Destroy() {
	log.Infof("gate %s 停止服务", g.opts.server.Addr())
	g.deregisterServiceInstance()

	g.stopNetworkServer()

	g.stopRPCServer()

	g.cancel()
}

// 启动网络服务器
func (g *Gate) startNetworkServer() {
	g.opts.server.OnConnect(g.handleConnect)
	g.opts.server.OnDisconnect(g.handleDisconnect)
	g.opts.server.OnReceive(g.handleReceive)

	if err := g.opts.server.Start(); err != nil {
		log.Fatalf("network server start failed: %v", err)
	}
}

// 停止网关服务器
func (g *Gate) stopNetworkServer() {
	if err := g.opts.server.Stop(); err != nil {
		log.Errorf("network server stop failed: %v", err)
	}
}

// 处理连接打开
func (g *Gate) handleConnect(conn network.Conn) {
	g.session.AddConn(conn)
	go func() {
		select {
		// todo be from config, from chan
		case <-time.After(time.Second * 6):
			if conn.UID() <= 0 {
				// 6秒 绑定上来，没进行登录操作的 则判定为攻击
				log.Errorf(" attack remoteip:%s,remoteAddr:%s", conn.RemoteIP(), conn.RemoteAddr())
				err := conn.Close()
				if err != nil {
					log.Errorf("connect not login err:%v", err)
				}
			} else {
				log.Debugf("auth has ,uid: %d", conn.UID())
			}
		}
	}()
}

// 处理断开连接
func (g *Gate) handleDisconnect(conn network.Conn) {
	g.session.RemConn(conn)

	if cid, uid := conn.ID(), conn.UID(); uid != 0 {
		ctx, cancel := context.WithTimeout(g.ctx, g.opts.timeout)
		_ = g.proxy.unbindGate(ctx, cid, uid)
		g.proxy.trigger(ctx, Disconnect, cid, uid)
		cancel()
	}

}

// 处理接收到的消息
func (g *Gate) handleReceive(conn network.Conn, data []byte) {
	cid, uid := conn.ID(), conn.UID()
	ctx, cancel := context.WithTimeout(g.ctx, g.opts.timeout)
	g.proxy.deliver(ctx, cid, uid, data)
	cancel()
}

// 启动RPC服务器
func (g *Gate) startRPCServer() {
	go func() {
		var err error
		err = g.rpcServer.Start(context.Background())
		if err != nil {
			log.Fatalf("rpc server create failed: %v", err)
		}
	}()
}

// 停止RPC服务器
func (g *Gate) stopRPCServer() {
	_ = g.rpcServer.Stop(context.Background())
}

// 注册服务实例
func (g *Gate) registerServiceInstance() {
	g.instance = &registry.ServiceInstance{
		ID:       g.opts.id,
		Name:     string(Gate),
		Kind:     Gate,
		Alias:    g.opts.name,
		State:    Work,
		Endpoint: g.rpcServer.Endpoint().String(),
	}

	ctx, cancel := context.WithTimeout(g.ctx, 10*time.Second)
	err := g.opts.registry.Register(ctx, g.instance)
	cancel()
	if err != nil {
		log.Fatalf("register dispatcher instance failed: %v", err)
	}
}

// 解注册服务实例
func (g *Gate) deregisterServiceInstance() {
	ctx, cancel := context.WithTimeout(g.ctx, 10*time.Second)
	defer cancel()
	err := g.opts.registry.Deregister(ctx, g.instance)
	if err != nil {
		log.Errorf("deregister dispatcher instance failed: %v", err)
	}
	log.Info("服务注册销毁成功")
}

func (g *Gate) debugPrint() {
	log.Debugf("gate server startup successful")
	log.Debugf("%s server listen on %s", g.opts.server.Protocol(), g.opts.server.Addr())
	log.Debugf("%s server listen on %s", g.rpcServer.Scheme(), g.rpcServer.Addr())
}

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
