package node

import (
	"context"

	"github.com/cr-mao/loric/cluster"
	"github.com/cr-mao/loric/packet"
	"github.com/cr-mao/loric/transport/grpc/code"
	"github.com/cr-mao/loric/transport/grpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// gate grpc serve 实现
type rpcService struct {
	pb.UnimplementedNodeServer
	provider *provider
}

// Trigger 触发事件
func (e *rpcService) Trigger(ctx context.Context, req *pb.TriggerRequest) (*pb.TriggerReply, error) {
	miss, err := e.provider.Trigger(ctx, &cluster.TriggerArgs{
		GID:   req.GID,
		CID:   req.CID,
		UID:   req.UID,
		Event: req.Event,
	})
	if err != nil {
		if miss {
			return nil, status.New(code.NotFoundSession, err.Error()).Err()
		} else {
			return nil, status.New(codes.Internal, err.Error()).Err()
		}
	}

	return &pb.TriggerReply{}, nil
}

// Deliver 投递消息
func (e *rpcService) Deliver(ctx context.Context, req *pb.DeliverRequest) (*pb.DeliverReply, error) {
	miss, err := e.provider.Deliver(ctx, &cluster.DeliverArgs{
		GID: req.GID,
		NID: req.NID,
		CID: req.CID,
		UID: req.UID,
		Message: &packet.Message{
			Seq:    req.Message.Seq,
			Route:  req.Message.Route,
			Buffer: req.Message.Buffer,
		},
	})
	if err != nil {
		if miss {
			return nil, status.New(code.NotFoundSession, err.Error()).Err()
		} else {
			return nil, status.New(codes.Internal, err.Error()).Err()
		}
	}
	return &pb.DeliverReply{}, nil
}
