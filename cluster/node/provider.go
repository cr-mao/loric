package node

import (
	"context"

	"github.com/cr-mao/loric/cluster"
)

type provider struct {
	node *Node
}

// Trigger 触发事件
func (p *provider) Trigger(ctx context.Context, args *cluster.TriggerArgs) (bool, error) {
	switch args.Event {
	case cluster.Reconnect:
		//  其实这里不会发生。目前gate 投递过来的事件 都是带用户id才投递
		if args.UID <= 0 {
			return false, ErrInvalidArgument
		}
		_, ok, err := p.node.proxy.AskNode(ctx, args.UID, p.node.opts.id)
		if err != nil {
			return false, err
		}
		if !ok {
			return true, ErrNotFoundUserSource
		}
	case cluster.Disconnect:
		if args.UID > 0 {
			_, ok, err := p.node.proxy.AskNode(ctx, args.UID, p.node.opts.id)
			if err != nil {
				return false, err
			}

			if !ok {
				return true, ErrNotFoundUserSource
			}
		}
	}

	handler, ok := p.node.events.events[args.Event]
	if !ok {
		return false, nil
	}

	evt := p.node.events.evtPool.Get().(*Event)
	evt.Event = args.Event
	evt.GID = args.GID
	evt.CID = args.CID
	evt.UID = args.UID
	defer p.node.events.evtPool.Put(evt)
	handler(evt)
	return false, nil
}

// Deliver 投递消息
func (p *provider) Deliver(ctx context.Context, args *cluster.DeliverArgs) (bool, error) {
	stateful, ok := p.node.router.CheckRouteStateful(args.Message.Route)
	if !ok {
		// 都给设置上，表示 404 处理消息 ,todo 要不要改成必须设置......
		if ok = p.node.router.HasDefaultRouteHandler(); !ok {
			return false, nil
		}
	}
	// 有状态都 用户id 必须传
	if stateful {
		if args.UID <= 0 {
			return false, ErrInvalidArgument
		}
		_, ok, err := p.node.proxy.AskNode(ctx, args.UID, p.node.opts.id)
		if err != nil {
			return false, err
		}
		if !ok {
			return true, ErrNotFoundUserSource
		}
	}
	p.node.router.deliver(args.GID, args.NID, args.CID, args.UID, args.Message.Seq, args.Message.Route, args.Message.Buffer)
	return false, nil
}
