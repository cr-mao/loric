package grpc

import (
	"sync"
	"sync/atomic"

	"github.com/cr-mao/loric/transport/grpc/resolver/direct"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
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
		return val.(*Pool).Get(), nil
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
	return pool.Get(), nil
}

type Pool struct {
	count uint64
	index uint64
	conns []*grpc.ClientConn
}

func newPool(count int, target string, opts ...grpc.DialOption) (*Pool, error) {
	p := &Pool{count: uint64(count), conns: make([]*grpc.ClientConn, count)}
	for i := 0; i < count; i++ {
		conn, err := grpc.Dial(target, opts...)
		if err != nil {
			return nil, err
		}
		p.conns[i] = conn
	}
	return p, nil
}

func (p *Pool) Get() *grpc.ClientConn {
	return p.conns[int(atomic.AddUint64(&p.index, 1)%p.count)]
}
