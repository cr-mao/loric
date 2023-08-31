package gate

import (
	"context"
	"sync"
	"time"

	"github.com/cr-mao/loric/cluster"
	"github.com/cr-mao/loric/errors"
	"github.com/cr-mao/loric/internal/dispatcher"
	"github.com/cr-mao/loric/internal/endpoint"
	"github.com/cr-mao/loric/locate"
	"github.com/cr-mao/loric/log"
	"github.com/cr-mao/loric/packet"
)

var ErrNotFoundUserSource = errors.New("not found user source")

type proxy struct {
	gate           *Gate    // 网关
	nodeSource     sync.Map //  用户在哪台 node
	nodeDispatcher *dispatcher.Dispatcher
}

func newProxy(gate *Gate) *proxy {
	return &proxy{
		gate:           gate,
		nodeDispatcher: dispatcher.NewDispatcher(dispatcher.RoundRobin),
	}
}

// 绑定用户与网关间的关系
func (p *proxy) bindGate(ctx context.Context, cid, uid int64) error {
	err := p.gate.opts.locator.Set(ctx, uid, cluster.Gate, p.gate.opts.id)
	if err != nil {
		return err
	}
	// 绑定 网关算重连， 登录的时候 ，node auth 后，则触发绑定网关的消息...
	p.trigger(ctx, cluster.Reconnect, cid, uid)
	return nil
}

// 解绑用户与网关间的关系
func (p *proxy) unbindGate(ctx context.Context, cid, uid int64) error {
	err := p.gate.opts.locator.Rem(ctx, uid, cluster.Gate, p.gate.opts.id)
	if err != nil {
		log.Errorf("user unbind failed, gid: %s, cid: %d, uid: %d, err: %v", p.gate.opts.id, cid, uid, err)
	}
	return err
}

// 触发事件
func (p *proxy) trigger(ctx context.Context, event int32, cid, uid int64) {
	if err := p.trigger2(ctx, &cluster.TriggerArgs{
		Event: event,
		CID:   cid,
		UID:   uid,
	}); err != nil {
		log.Warnf("trigger event failed, gid: %s, cid: %d, uid: %d, event: [%s], err: %v", p.gate.opts.id, cid, uid, cluster.EventNames[event], err)
	}
}

// Trigger 触发事件
func (p *proxy) trigger2(ctx context.Context, args *cluster.TriggerArgs) error {
	// switch 这块其实都不会走，先留着
	switch args.Event {
	// 这里不会走到
	case cluster.Connect:
		return p.doTrigger(ctx, args)
	case cluster.Disconnect:
		// 限定了 必须有用户id才触发
		if args.UID == 0 {
			return p.doTrigger(ctx, args)
		}
		// 绑定网关
	case cluster.Reconnect:
		if args.UID == 0 {
			return ErrInvalidArgument
		}
	}

	var (
		err       error
		nid       string
		prev      string
		miss      bool
		client    *NodeGrpcClient
		ep        *endpoint.Endpoint
		arguments = &cluster.TriggerArgs{
			Event: args.Event,
			GID:   p.gate.opts.id,
			CID:   args.CID,
			UID:   args.UID,
		}
	)

	for i := 0; i < 2; i++ {
		if nid, err = p.LocateNode(ctx, args.UID); err != nil {
			//if args.Event == cluster.Disconnect && err == ErrNotFoundUserSource {
			//	//todo 这个case 好想没用，用户都不知道再哪台机器上....
			//	return p.doTrigger(ctx, args)
			//}
			return err
		}
		if nid == prev {
			return err
		}
		prev = nid

		if ep, err = p.nodeDispatcher.FindEndpoint(nid); err != nil {
			// 这个不太可能出现，
			if args.Event == cluster.Disconnect && err == dispatcher.ErrNotFoundEndpoint {
				return p.doTrigger(ctx, args)
			}
			return err
		}
		client, err = p.gate.opts.transporter.NewNodeClient(ep)
		if err != nil {
			return err
		}

		miss, err = client.Trigger(ctx, arguments)
		if miss {
			p.nodeSource.Delete(args.UID)
			continue
		}

		break
	}

	return err
}

// LocateNode 定位用户所在节点
func (p *proxy) LocateNode(ctx context.Context, uid int64) (string, error) {
	if val, ok := p.nodeSource.Load(uid); ok {
		if nid := val.(string); nid != "" {
			return nid, nil
		}
	}

	nid, err := p.gate.opts.locator.Get(ctx, uid, cluster.Node)

	if err != nil {
		return "", err
	}

	if nid == "" {
		p.nodeSource.Delete(uid)
		return "", ErrNotFoundUserSource
	}
	p.nodeSource.Store(uid, nid)
	return nid, nil
}

// 触发事件
func (p *proxy) doTrigger(ctx context.Context, args *cluster.TriggerArgs) error {
	event, err := p.nodeDispatcher.FindEvent(args.Event)
	if err != nil {
		if err == dispatcher.ErrNotFoundEvent {
			return nil
		}

		return err
	}

	ep, err := event.FindEndpoint()
	if err != nil {
		if err == dispatcher.ErrNotFoundEndpoint {
			return nil
		}

		return err
	}

	client, err := p.gate.opts.transporter.NewNodeClient(ep)
	if err != nil {
		return err
	}

	_, err = client.Trigger(ctx, &cluster.TriggerArgs{
		Event: args.Event,
		GID:   p.gate.opts.id,
		CID:   args.CID,
		UID:   args.UID,
	})
	return err
}

// 投递消息给节点
func (p *proxy) deliver(ctx context.Context, cid, uid int64, data []byte) {
	message, err := packet.Unpack(data)
	if err != nil {
		log.Errorf("unpack data to struct failed: %v", err)
		return
	}
	arguments := &cluster.DeliverArgs{
		GID:     p.gate.opts.id,
		NID:     "",
		CID:     cid,
		UID:     uid,
		Message: message,
	}

	_, err = p.doNodeRPC(ctx, message.Route, uid, func(ctx context.Context, client *NodeGrpcClient) (bool, interface{}, error) {
		miss, errDelver := client.Deliver(ctx, arguments)
		return miss, nil, errDelver
	})
	if err != nil {
		log.Errorf("deliver message failed: %v", err)
		return
	}
}

// 执行节点RPC调用, todo 简化， 这代码不够好阅读
func (p *proxy) doNodeRPC(ctx context.Context, routeID int32, uid int64, fn func(ctx context.Context, client *NodeGrpcClient) (bool, interface{}, error)) (interface{}, error) {
	var (
		err       error
		nid       string
		prev      string
		client    *NodeGrpcClient
		route     *dispatcher.Route
		ep        *endpoint.Endpoint
		continued bool
		reply     interface{}
	)

	if route, err = p.nodeDispatcher.FindRoute(routeID); err != nil {
		return nil, err
	}

	for i := 0; i < 2; i++ {
		if route.Stateful() {
			if nid, err = p.LocateNode(ctx, uid); err != nil {
				return nil, err
			}
			if nid == prev {
				return reply, err
			}
			prev = nid
		}

		ep, err = route.FindEndpoint(nid)
		if err != nil {
			return nil, err
		}

		client, err = p.gate.opts.transporter.NewNodeClient(ep)
		if err != nil {
			return nil, err
		}

		continued, reply, err = fn(ctx, client)
		if continued {
			p.nodeSource.Delete(uid)
			continue
		}

		break
	}

	return reply, err
}

// 启动监听
func (p *proxy) watch(ctx context.Context) {
	p.watchUserLocate(ctx)
	p.watchServiceInstance(ctx)
}

// 监听服务实例
func (p *proxy) watchServiceInstance(ctx context.Context) {
	rctx, rcancel := context.WithTimeout(ctx, 10*time.Second)
	watcher, err := p.gate.opts.registry.Watch(rctx, cluster.Node)
	rcancel()
	if err != nil {
		log.Fatalf("the dispatcher instance watch failed: %v", err)
	}
	go func() {
		defer watcher.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// exec watch
			}
			services, err := watcher.Next()
			if err != nil {
				continue
			}
			p.nodeDispatcher.ReplaceServices(services...)
		}
	}()
}

// WatchUserLocate 监听用户定位 ,只监听 用户node 节点的变化
func (p *proxy) watchUserLocate(ctx context.Context) {
	rctx, rcancel := context.WithTimeout(ctx, 10*time.Second)
	watcher, err := p.gate.opts.locator.Watch(rctx, cluster.Node)
	rcancel()
	if err != nil {
		log.Fatalf("user locate event watch failed: %v", err)
	}
	go func() {
		defer watcher.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// exec watch
			}
			events, err := watcher.Next()
			if err != nil {
				continue
			}
			for _, event := range events {
				var source *sync.Map
				switch event.InsKind {
				case cluster.Node:
					source = &p.nodeSource
				}
				if source == nil {
					continue
				}
				switch event.Type {
				case locate.SetLocation:
					source.Store(event.UID, event.InsID)
				case locate.RemLocation:
					source.Delete(event.UID)
				}
			}
		}
	}()
}
