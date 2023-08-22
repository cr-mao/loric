/**
User: cr-mao
Date: 2023/8/19 10:38
Email: crmao@qq.com
Desc: grpc_server.go
*/
package gate

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cr-mao/loric/packet"
	"github.com/cr-mao/loric/session"
	"github.com/cr-mao/loric/transport/grpc/code"
	"github.com/cr-mao/loric/transport/grpc/pb"
)

// gate grpc serve 实现
type rpcService struct {
	pb.UnimplementedGateServer
	provider *provider
}

// Bind 将用户与当前网关进行绑定
func (e *rpcService) Bind(ctx context.Context, req *pb.BindRequest) (*pb.BindReply, error) {
	err := e.provider.Bind(ctx, req.CID, req.UID)
	if err != nil {
		switch err {
		case session.ErrNotFoundSession:
			return nil, status.New(code.NotFoundSession, err.Error()).Err()
		case session.ErrInvalidSessionKind:
			return nil, status.New(codes.InvalidArgument, err.Error()).Err()
		default:
			return nil, status.New(codes.Internal, err.Error()).Err()
		}
	}
	return &pb.BindReply{}, nil
}

// Unbind 将用户与当前网关进行解绑
func (e *rpcService) Unbind(ctx context.Context, req *pb.UnbindRequest) (*pb.UnbindReply, error) {
	err := e.provider.Unbind(ctx, req.UID)
	if err != nil {
		switch err {
		case session.ErrNotFoundSession:
			return nil, status.New(code.NotFoundSession, err.Error()).Err()
		case session.ErrInvalidSessionKind:
			return nil, status.New(codes.InvalidArgument, err.Error()).Err()
		default:
			return nil, status.New(codes.Internal, err.Error()).Err()
		}
	}

	return &pb.UnbindReply{}, nil
}

// GetIP 获取客户端IP地址
func (e *rpcService) GetIP(ctx context.Context, req *pb.GetIPRequest) (*pb.GetIPReply, error) {
	ip, err := e.provider.GetIP(ctx, session.Kind(req.Kind), req.Target)
	if err != nil {
		switch err {
		case session.ErrNotFoundSession:
			return nil, status.New(code.NotFoundSession, err.Error()).Err()
		case session.ErrInvalidSessionKind:
			return nil, status.New(codes.InvalidArgument, err.Error()).Err()
		default:
			return nil, status.New(codes.Internal, err.Error()).Err()
		}
	}

	return &pb.GetIPReply{IP: ip}, nil
}

// Push 推送消息给连接
func (e *rpcService) Push(ctx context.Context, req *pb.PushRequest) (*pb.PushReply, error) {
	err := e.provider.Push(ctx, session.Kind(req.Kind), req.Target, &packet.Message{
		Seq:    req.Message.Seq,
		Route:  req.Message.Route,
		Buffer: req.Message.Buffer,
	})
	if err != nil {
		switch err {
		case session.ErrNotFoundSession:
			return nil, status.New(code.NotFoundSession, err.Error()).Err()
		case session.ErrInvalidSessionKind:
			return nil, status.New(codes.InvalidArgument, err.Error()).Err()
		default:
			return nil, status.New(codes.Internal, err.Error()).Err()
		}
	}

	return &pb.PushReply{}, nil
}

// Multicast 推送组播消息
func (e *rpcService) Multicast(ctx context.Context, req *pb.MulticastRequest) (*pb.MulticastReply, error) {
	total, err := e.provider.Multicast(ctx, session.Kind(req.Kind), req.Targets, &packet.Message{
		Seq:    req.Message.Seq,
		Route:  req.Message.Route,
		Buffer: req.Message.Buffer,
	})
	if err != nil {
		switch err {
		case session.ErrInvalidSessionKind:
			return nil, status.New(codes.InvalidArgument, err.Error()).Err()
		default:
			return nil, status.New(codes.Internal, err.Error()).Err()
		}
	}

	return &pb.MulticastReply{Total: total}, nil
}

// Broadcast 推送广播消息
func (e *rpcService) Broadcast(ctx context.Context, req *pb.BroadcastRequest) (*pb.BroadcastReply, error) {
	total, err := e.provider.Broadcast(ctx, session.Kind(req.Kind), &packet.Message{
		Seq:    req.Message.Seq,
		Route:  req.Message.Route,
		Buffer: req.Message.Buffer,
	})
	if err != nil {
		switch err {
		case session.ErrInvalidSessionKind:
			return nil, status.New(codes.InvalidArgument, err.Error()).Err()
		default:
			return nil, status.New(codes.Internal, err.Error()).Err()
		}
	}

	return &pb.BroadcastReply{Total: total}, nil
}

// Stat 统计会话总数
func (e *rpcService) Stat(ctx context.Context, req *pb.StatRequest) (*pb.StatReply, error) {
	total, err := e.provider.Stat(ctx, session.Kind(req.Kind))
	if err != nil {
		switch err {
		case session.ErrInvalidSessionKind:
			return nil, status.New(codes.InvalidArgument, err.Error()).Err()
		default:
			return nil, status.New(codes.Internal, err.Error()).Err()
		}
	}

	return &pb.StatReply{Total: total}, nil
}

// Disconnect 断开连接
func (e *rpcService) Disconnect(ctx context.Context, req *pb.DisconnectRequest) (*pb.DisconnectReply, error) {
	err := e.provider.Disconnect(ctx, session.Kind(req.Kind), req.Target, req.IsForce)
	if err != nil {
		switch err {
		case session.ErrNotFoundSession:
			return nil, status.New(code.NotFoundSession, err.Error()).Err()
		case session.ErrInvalidSessionKind:
			return nil, status.New(codes.InvalidArgument, err.Error()).Err()
		default:
			return nil, status.New(codes.Internal, err.Error()).Err()
		}
	}
	return &pb.DisconnectReply{}, nil
}
