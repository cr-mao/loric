/**
User: cr-mao
Date: 2023/8/17 14:27
Email: crmao@qq.com
Desc: server.go
*/
package grpc

import (
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/cr-mao/loric/internal/endpoint"
	"github.com/cr-mao/loric/internal/netlib"
	"github.com/cr-mao/loric/log"
	"github.com/cr-mao/loric/transport"
)

const scheme = "grpc"

var _ transport.Server = (*Server)(nil)

// Server is a gRPC server wrapper.
type Server struct {
	*grpc.Server
	listenAddr string                        // 监听地址
	exposeAddr string                        // 内网地址
	unaryInts  []grpc.UnaryServerInterceptor //一次元拦截器
	grpcOpts   []grpc.ServerOption           //
	health     *health.Server                // 健康检测server
	endpoint   *endpoint.Endpoint
}

// NewServer creates a gRPC server by options.
func NewServer(addr string, opts ...ServerOption) (*Server, error) {
	listenAddr, exposeAddr, err := netlib.ParseAddr(addr)
	if err != nil {
		return nil, err
	}
	srv := &Server{
		health:     health.NewServer(),
		listenAddr: listenAddr,
		exposeAddr: exposeAddr,
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
	srv.endpoint = endpoint.NewEndpoint(scheme, exposeAddr, false)
	return srv, nil
}

// Endpoint return a real address to registry endpoint.
// examples:
//	grpc://127.0.0.1:9000?is_secure=false,
func (s *Server) Endpoint() *endpoint.Endpoint {
	return s.endpoint
}

// Addr 监听地址
func (s *Server) Addr() string {
	return s.listenAddr
}

// Scheme 协议
func (s *Server) Scheme() string {
	return scheme
}

// Start start the gRPC server.
func (s *Server) Start() error {
	addr, err := net.ResolveTCPAddr("tcp", s.listenAddr)
	if err != nil {
		return err
	}

	lis, err := net.Listen(addr.Network(), addr.String())
	if err != nil {
		return err
	}
	//设置serving 状态
	s.health.Resume()
	return s.Server.Serve(lis)
}

// Stop stop the gRPC server.
func (s *Server) Stop() error {
	s.health.Shutdown()
	s.Server.GracefulStop()
	log.Info("[gRPC] server stopping")
	return nil
}

// ServerOption is gRPC server option.
type ServerOption func(o *Server)

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
