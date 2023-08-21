/**
User: cr-mao
Date: 2023/8/18 15:21
Email: crmao@qq.com
Desc: node_client.go
*/
package gate

import (
	"context"

	"github.com/cr-mao/loric/cluster"
	"github.com/cr-mao/loric/transport/grpc/code"
	"github.com/cr-mao/loric/transport/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/status"
)

var _ cluster.NodeClient = (*NodeGrpcClient)(nil)

type NodeGrpcClient struct {
	client pb.NodeClient
}

func NewNodeClient(cc *grpc.ClientConn) *NodeGrpcClient {
	return &NodeGrpcClient{client: pb.NewNodeClient(cc)}
}

// Trigger 触发事件
func (c *NodeGrpcClient) Trigger(ctx context.Context, args *cluster.TriggerArgs) (miss bool, err error) {
	_, err = c.client.Trigger(ctx, &pb.TriggerRequest{
		Event: args.Event,
		GID:   args.GID,
		CID:   args.CID,
		UID:   args.UID,
	})
	miss = status.Code(err) == code.NotFoundSession
	return
}

// Deliver 投递消息
func (c *NodeGrpcClient) Deliver(ctx context.Context, args *cluster.DeliverArgs) (miss bool, err error) {
	_, err = c.client.Deliver(ctx, &pb.DeliverRequest{
		GID: args.GID,
		NID: args.NID,
		CID: args.CID,
		UID: args.UID,
		Message: &pb.Message{
			Seq:    args.Message.Seq,
			Route:  args.Message.Route,
			Buffer: args.Message.Buffer,
		},
	}, grpc.UseCompressor(gzip.Name))
	miss = status.Code(err) == code.NotFoundSession
	return
}
