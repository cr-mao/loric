/**
User: cr-mao
Date: 2023/8/19 03:39
Email: crmao@qq.com
Desc: transport.go
*/
package cluster

import (
	"context"
	"github.com/cr-mao/loric/packet"
	"github.com/cr-mao/loric/session"
)

type NodeClient interface {
	// Trigger 触发事件
	Trigger(ctx context.Context, args *TriggerArgs) (miss bool, err error)
	// Deliver 投递消息
	Deliver(ctx context.Context, args *DeliverArgs) (miss bool, err error)
}

type GateClient interface {
	// Bind 绑定用户与连接
	Bind(ctx context.Context, cid, uid int64) (miss bool, err error)
	// Unbind 解绑用户与连接
	Unbind(ctx context.Context, uid int64) (miss bool, err error)
	// GetIP 获取客户端IP
	GetIP(ctx context.Context, kind session.Kind, target int64) (ip string, miss bool, err error)
	// Push 推送消息
	Push(ctx context.Context, kind session.Kind, target int64, message *packet.Message) (miss bool, err error)
	// Multicast 推送组播消息
	Multicast(ctx context.Context, kind session.Kind, targets []int64, message *packet.Message) (total int64, err error)
	// Broadcast 推送广播消息
	Broadcast(ctx context.Context, kind session.Kind, message *packet.Message) (total int64, err error)
	// Stat 统计会话总数
	Stat(ctx context.Context, kind session.Kind) (total int64, err error)
	// Disconnect 断开连接
	Disconnect(ctx context.Context, kind session.Kind, target int64, isForce bool) (miss bool, err error)
}

type GateProvider interface {
	// Bind 绑定用户与网关间的关系
	Bind(ctx context.Context, cid, uid int64) error
	// Unbind 解绑用户与网关间的关系
	Unbind(ctx context.Context, uid int64) error
	// GetIP 获取客户端IP地址
	GetIP(ctx context.Context, kind session.Kind, target int64) (ip string, err error)
	// Push 发送消息（异步）
	Push(ctx context.Context, kind session.Kind, target int64, message *packet.Message) error
	// Multicast 推送组播消息（异步）
	Multicast(ctx context.Context, kind session.Kind, targets []int64, message *packet.Message) (total int64, err error)
	// Broadcast 推送广播消息（异步）
	Broadcast(ctx context.Context, kind session.Kind, message *packet.Message) (total int64, err error)
	// Stat 统计会话总数
	Stat(ctx context.Context, kind session.Kind) (total int64, err error)
	// Disconnect 断开连接
	Disconnect(ctx context.Context, kind session.Kind, target int64, isForce bool) error
}

type NodeProvider interface {
	// Trigger 触发事件
	Trigger(ctx context.Context, args *TriggerArgs) (miss bool, err error)
	// Deliver 投递消息
	Deliver(ctx context.Context, args *DeliverArgs) (miss bool, err error)
}

type DeliverArgs struct {
	GID     string
	NID     string
	CID     int64
	UID     int64
	Message *packet.Message
}

type TriggerArgs struct {
	Event Event
	GID   string
	CID   int64
	UID   int64
}
