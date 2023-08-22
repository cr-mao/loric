/**
User: cr-mao
Date: 2023/8/21 15:03
Email: crmao@qq.com
Desc: gate_client.go
*/
package node

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/status"

	"github.com/cr-mao/loric/cluster"
	"github.com/cr-mao/loric/packet"
	"github.com/cr-mao/loric/session"
	"github.com/cr-mao/loric/transport/grpc/code"
	"github.com/cr-mao/loric/transport/grpc/pb"
)

var _ cluster.GateClient = (*GateGrpcClient)(nil)

type GateGrpcClient struct {
	client pb.GateClient
}

func NewGateClient(cc *grpc.ClientConn) *GateGrpcClient {
	return &GateGrpcClient{client: pb.NewGateClient(cc)}
}

// Bind 绑定用户与连接
func (c *GateGrpcClient) Bind(ctx context.Context, cid, uid int64) (miss bool, err error) {
	_, err = c.client.Bind(ctx, &pb.BindRequest{
		CID: cid,
		UID: uid,
	})
	miss = status.Code(err) == code.NotFoundSession
	return
}

// Unbind 解绑用户与连接
func (c *GateGrpcClient) Unbind(ctx context.Context, uid int64) (miss bool, err error) {
	_, err = c.client.Unbind(ctx, &pb.UnbindRequest{
		UID: uid,
	})

	miss = status.Code(err) == code.NotFoundSession

	return
}

// GetIP 获取客户端IP
func (c *GateGrpcClient) GetIP(ctx context.Context, kind session.Kind, target int64) (ip string, miss bool, err error) {
	reply, err := c.client.GetIP(ctx, &pb.GetIPRequest{
		Kind:   int32(kind),
		Target: target,
	})
	if err != nil {
		miss = status.Code(err) == code.NotFoundSession
		return
	}

	ip = reply.IP

	return
}

// Push 推送消息
func (c *GateGrpcClient) Push(ctx context.Context, kind session.Kind, target int64, message *packet.Message) (miss bool, err error) {
	_, err = c.client.Push(ctx, &pb.PushRequest{
		Kind:   int32(kind),
		Target: target,
		Message: &pb.Message{
			Seq:    message.Seq,
			Route:  message.Route,
			Buffer: message.Buffer,
		},
	}, grpc.UseCompressor(gzip.Name))

	miss = status.Code(err) == code.NotFoundSession

	return
}

// Multicast 推送组播消息
func (c *GateGrpcClient) Multicast(ctx context.Context, kind session.Kind, targets []int64, message *packet.Message) (int64, error) {
	reply, err := c.client.Multicast(ctx, &pb.MulticastRequest{
		Kind:    int32(kind),
		Targets: targets,
		Message: &pb.Message{
			Seq:    message.Seq,
			Route:  message.Route,
			Buffer: message.Buffer,
		},
	}, grpc.UseCompressor(gzip.Name))
	if err != nil {
		return 0, err
	}

	return reply.Total, nil
}

// Broadcast 推送广播消息
func (c *GateGrpcClient) Broadcast(ctx context.Context, kind session.Kind, message *packet.Message) (int64, error) {
	reply, err := c.client.Broadcast(ctx, &pb.BroadcastRequest{
		Kind: int32(kind),
		Message: &pb.Message{
			Seq:    message.Seq,
			Route:  message.Route,
			Buffer: message.Buffer,
		},
	}, grpc.UseCompressor(gzip.Name))
	if err != nil {
		return 0, err
	}

	return reply.Total, nil
}

// Stat 统计会话总数
func (c *GateGrpcClient) Stat(ctx context.Context, kind session.Kind) (int64, error) {
	reply, err := c.client.Stat(ctx, &pb.StatRequest{
		Kind: int32(kind),
	})
	if err != nil {
		return 0, err
	}
	return reply.Total, nil
}

// Disconnect 断开连接
func (c *GateGrpcClient) Disconnect(ctx context.Context, kind session.Kind, target int64, isForce bool) (miss bool, err error) {
	_, err = c.client.Disconnect(ctx, &pb.DisconnectRequest{
		Kind:    int32(kind),
		Target:  target,
		IsForce: isForce,
	})
	miss = status.Code(err) == code.NotFoundSession
	return
}
