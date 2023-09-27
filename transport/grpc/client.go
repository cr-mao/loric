package grpc

import (
	"errors"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"

	"github.com/cr-mao/loric/transport/grpc/resolver/direct"
)

type ClientBuilder struct {
	opts     *ClientOptions
	dialOpts []grpc.DialOption
	pools    sync.Map
}

type ClientOptions struct {
	PoolSize int
	DialOpts []grpc.DialOption
}

func NewClientBuilder(opts *ClientOptions) *ClientBuilder {
	b := &ClientBuilder{opts: opts}
	// 暂时统一不加密传输
	var creds = insecure.NewCredentials()
	resolvers := make([]resolver.Builder, 0, 1)
	// 暂时统一direct直连模式
	resolvers = append(resolvers, direct.NewBuilder())
	b.dialOpts = make([]grpc.DialOption, 0, len(opts.DialOpts)+2)
	b.dialOpts = append(b.dialOpts, grpc.WithTransportCredentials(creds))
	b.dialOpts = append(b.dialOpts, grpc.WithResolvers(resolvers...))
	return b
}

// Build 构建连接
func (b *ClientBuilder) GetConn(target string) (*grpc.ClientConn, error) {
	val, ok := b.pools.Load(target)
	if ok {
		return val.(*Pool).Get()
	}
	size := b.opts.PoolSize
	if size <= 0 {
		size = 10
	}
	pool, err := newPool(size, target, b.dialOpts...)
	if err != nil {
		return nil, err
	}
	b.pools.Store(target, pool)
	return pool.Get()
}

type Pool struct {
	target string
	count  uint64
	index  uint64
	conns  []*grpc.ClientConn
	sync.Mutex
	opts []grpc.DialOption
}

func newPool(count int, target string, opts ...grpc.DialOption) (*Pool, error) {
	p := &Pool{
		count:  uint64(count),
		conns:  make([]*grpc.ClientConn, count),
		target: target,
		opts:   opts,
	}
	for i := 0; i < count; i++ {
		conn, err := grpc.Dial(target, opts...)
		if err != nil {
			return nil, err
		}
		p.conns[i] = conn
	}
	return p, nil
}

func (p *Pool) Get() (*grpc.ClientConn, error) {
	idx := int(atomic.AddUint64(&p.index, 1) % p.count)
	conn := p.conns[idx]
	if conn != nil && p.checkState(conn) == nil {
		return conn, nil
	}

	p.Lock()
	// gc old conn
	if conn != nil {
		conn.Close()
	}
	defer p.Unlock()
	// double check, already inited
	conn = p.conns[idx]
	if conn != nil && p.checkState(conn) == nil {
		return conn, nil
	}
	conn, err := grpc.Dial(p.target, p.opts...)
	if err != nil {
		return nil, err
	}
	p.conns[idx] = conn
	return conn, nil
}

var ErrConnShutdown = errors.New("grpc conn shutdown")

func (p *Pool) checkState(conn *grpc.ClientConn) error {
	state := conn.GetState()
	switch state {
	case connectivity.TransientFailure, connectivity.Shutdown:
		return ErrConnShutdown
	}
	return nil
}
