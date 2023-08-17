/**
User: cr-mao
Date: 2023/8/17 14:27
Email: crmao@qq.com
Desc: server.go
*/
package grpc

import (
	"context"
	"net"
	"net/url"
	"time"

	"github.com/cr-mao/loric/internal/endpoint"
	"github.com/cr-mao/loric/internal/host"
	"github.com/cr-mao/loric/log"
	"github.com/cr-mao/loric/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var _ transport.Endpointer = (*Server)(nil)
var _ transport.Server = (*Server)(nil)

// Server is a gRPC server wrapper.
type Server struct {
	*grpc.Server
	baseCtx   context.Context
	address   string                        //地址
	unaryInts []grpc.UnaryServerInterceptor //一次元拦截器
	grpcOpts  []grpc.ServerOption           //
	lis       net.Listener
	timeout   time.Duration
	health    *health.Server // 健康检测server
	endpoint  *url.URL       // url
	err       error
}

// NewServer creates a gRPC server by options.
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		baseCtx: context.Background(),
		address: ":0",
		timeout: 2 * time.Second,
		health:  health.NewServer(),
	}
	for _, o := range opts {
		o(srv)
	}
	unaryInts := []grpc.UnaryServerInterceptor{}
	if len(srv.unaryInts) > 0 {
		unaryInts = append(unaryInts, srv.unaryInts...)
	}

	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(unaryInts...),
	}
	//用户传的ServerOption
	if len(srv.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, srv.grpcOpts...)
	}
	srv.Server = grpc.NewServer(grpcOpts...)
	grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)
	return srv
}

// Endpoint return a real address to registry endpoint.
// examples:
//	grpc://127.0.0.1:9000?isSecure=false
func (s *Server) Endpoint() (*url.URL, error) {
	if err := s.listenAndEndpoint(); err != nil {
		return nil, s.err
	}
	return s.endpoint, nil
}

// Start start the gRPC server.
func (s *Server) Start(ctx context.Context) error {
	if err := s.listenAndEndpoint(); err != nil {
		return s.err
	}
	s.baseCtx = ctx
	log.Infof("[gRPC] server listening on: %s", s.lis.Addr().String())
	//设置serving 状态
	s.health.Resume()
	return s.Serve(s.lis)
}

// Stop stop the gRPC server.
func (s *Server) Stop(_ context.Context) error {
	s.health.Shutdown()
	s.GracefulStop()
	log.Info("[gRPC] server stopping")
	return nil
}

func (s *Server) listenAndEndpoint() error {
	if s.lis == nil {
		lis, err := net.Listen("tcp", s.address)
		if err != nil {
			s.err = err
			return err
		}
		s.lis = lis
	}
	if s.endpoint == nil {
		addr, err := host.Extract(s.address, s.lis)
		if err != nil {
			s.err = err
			return err
		}
		s.endpoint = endpoint.NewEndpoint(endpoint.Scheme("grpc", false), addr)
	}
	return s.err
}

// ServerOption is gRPC server option.
type ServerOption func(o *Server)

func WithAddress(address string) ServerOption {
	return func(o *Server) {
		o.address = address
	}
}

// Listener with server lis
func WithListener(lis net.Listener) ServerOption {
	return func(s *Server) {
		s.lis = lis
	}
}

// WithUnaryInterceptor returns a ServerOption that sets the UnaryServerInterceptor for the server.
func WithUnaryInterceptor(in ...grpc.UnaryServerInterceptor) ServerOption {
	return func(s *Server) {
		s.unaryInts = in
	}
}

// WithGrpcOpts with grpc options.
func WithGrpcOpts(opts ...grpc.ServerOption) ServerOption {
	return func(s *Server) {
		s.grpcOpts = opts
	}
}
