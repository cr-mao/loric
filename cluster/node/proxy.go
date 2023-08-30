package node

import (
	"context"
	"github.com/cr-mao/loric/sugar/slice"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/cr-mao/loric/cluster"
	"github.com/cr-mao/loric/internal/dispatcher"
	"github.com/cr-mao/loric/internal/endpoint"
	"github.com/cr-mao/loric/locate"
	"github.com/cr-mao/loric/log"
	"github.com/cr-mao/loric/packet"
	"github.com/cr-mao/loric/registry"
	"github.com/cr-mao/loric/session"
)

type Proxy struct {
	node           *Node    // 节点
	nodeSource     sync.Map //  用户在哪台 node
	gateSource     sync.Map //  用户来源网关
	nodeDispatcher *dispatcher.Dispatcher
	gateDispatcher *dispatcher.Dispatcher
}

func newProxy(node *Node) *Proxy {
	return &Proxy{
		node: node,
		//todo  先写死策略
		nodeDispatcher: dispatcher.NewDispatcher(dispatcher.RoundRobin),
		gateDispatcher: dispatcher.NewDispatcher(dispatcher.RoundRobin),
	}
}

// GetNodeID 获取当前节点ID
func (p *Proxy) GetNodeID() string {
	return p.node.opts.id
}

// GetNodeName 获取当前节点名称
func (p *Proxy) GetNodeName() string {
	return p.node.opts.name
}

// GetNodeState 获取当前节点状态
func (p *Proxy) GetNodeState() int32 {
	return p.node.getState()
}

// SetNodeState 设置当前节点状态
func (p *Proxy) SetNodeState(state int32) {
	p.node.setState(state)
}

// Router 路由器
func (p *Proxy) Router() *Router {
	return p.node.router
}

// Events 事件分发器
func (p *Proxy) Events() *Events {
	return p.node.events
}

// BindGate 绑定网关
func (p *Proxy) BindGate(ctx context.Context, uid int64, gid string, cid int64) error {
	client, err := p.getGateClientByGID(gid)
	if err != nil {
		return err
	}
	_, err = client.Bind(ctx, cid, uid)
	if err != nil {
		return err
	}
	p.gateSource.Store(uid, gid)
	return nil
}

// 根据实例ID获取网关客户端
func (p *Proxy) getGateClientByGID(gid string) (*GateGrpcClient, error) {
	if gid == "" {
		return nil, ErrInvalidGID
	}
	ep, err := p.gateDispatcher.FindEndpoint(gid)
	if err != nil {
		return nil, err
	}
	return p.node.opts.transporter.NewGateClient(ep)
}

// 根据实例ID获取节点客户端
func (p *Proxy) getNodeClientByNID(nid string) (*NodeGrpcClient, error) {
	if nid == "" {
		return nil, ErrInvalidNID
	}

	ep, err := p.nodeDispatcher.FindEndpoint(nid)
	if err != nil {
		return nil, err
	}

	return p.node.opts.transporter.NewNodeClient(ep)
}

// UnbindGate 解绑网关
func (p *Proxy) UnbindGate(ctx context.Context, uid int64) error {
	_, err := p.doGateRPC(ctx, uid, func(client *GateGrpcClient) (bool, interface{}, error) {
		miss, err := client.Unbind(ctx, uid)
		return miss, nil, err
	})
	if err != nil {
		return err
	}

	p.gateSource.Delete(uid)

	return nil
}

// 执行网关RPC调用
func (p *Proxy) doGateRPC(ctx context.Context, uid int64, fn func(client *GateGrpcClient) (bool, interface{}, error)) (interface{}, error) {
	var (
		err       error
		gid       string
		prev      string
		client    *GateGrpcClient
		continued bool
		reply     interface{}
	)

	for i := 0; i < 2; i++ {
		if gid, err = p.LocateGate(ctx, uid); err != nil {
			return nil, err
		}
		if gid == prev {
			return reply, err
		}
		prev = gid
		client, err = p.getGateClientByGID(gid)
		if err != nil {
			return nil, err
		}
		continued, reply, err = fn(client)
		if continued {
			p.gateSource.Delete(uid)
			continue
		}
		break
	}

	return reply, err
}

// BindNode 绑定节点
// 单个用户只能被绑定到某一台节点服务器上，多次绑定会直接覆盖上次绑定
// 绑定操作会通过发布订阅方式同步到网关服务器和其他相关节点服务器上
// NID 为需要绑定的节点ID，默认绑定到当前节点上
func (p *Proxy) BindNode(ctx context.Context, uid int64, nid ...string) error {
	var bindNid = p.node.opts.id
	if len(nid) > 0 && nid[0] != "" {
		bindNid = nid[0]
	}
	err := p.node.opts.locator.Set(ctx, uid, cluster.Node, bindNid)
	if err != nil {
		return err
	}
	p.nodeSource.Store(uid, bindNid)
	return nil
}

// UnbindNode 解绑节点
// 解绑时会对解绑节点ID进行校验，不匹配则解绑失败
// 解绑操作会通过发布订阅方式同步到网关服务器和其他相关节点服务器上
// NID 为需要解绑的节点ID，默认解绑当前节点
func (p *Proxy) UnbindNode(ctx context.Context, uid int64, nid ...string) error {
	var unbindNid = p.node.opts.id
	if len(nid) > 0 && nid[0] != "" {
		unbindNid = nid[0]
	}
	err := p.node.opts.locator.Rem(ctx, uid, cluster.Node, unbindNid)
	if err != nil {
		return err
	}
	p.nodeSource.Delete(uid)
	return nil
}

// LocateGate 定位用户所在网关
func (p *Proxy) LocateGate(ctx context.Context, uid int64) (string, error) {
	if val, ok := p.gateSource.Load(uid); ok {
		if gid := val.(string); gid != "" {
			return gid, nil
		}
	}
	gid, err := p.node.opts.locator.Get(ctx, uid, cluster.Gate)
	if err != nil {
		return "", err
	}
	if gid == "" {
		p.gateSource.Delete(uid)

		return "", ErrNotFoundUserSource
	}
	p.gateSource.Store(uid, gid)
	return gid, nil
}

// AskGate 检测用户是否在给定的网关上
func (p *Proxy) AskGate(ctx context.Context, uid int64, gid string) (string, bool, error) {
	if val, ok := p.gateSource.Load(uid); ok {
		if val.(string) == gid {
			return gid, true, nil
		}
	}
	insID, err := p.node.opts.locator.Get(ctx, uid, cluster.Gate)
	if err != nil {
		return "", false, err
	}
	if insID == "" {
		p.gateSource.Delete(uid)
		return "", false, ErrNotFoundUserSource
	}
	p.gateSource.Store(uid, insID)
	return insID, insID == gid, nil
}

// LocateNode 定位用户所在节点
func (p *Proxy) LocateNode(ctx context.Context, uid int64) (string, error) {
	if val, ok := p.nodeSource.Load(uid); ok {
		if nid := val.(string); nid != "" {
			return nid, nil
		}
	}
	nid, err := p.node.opts.locator.Get(ctx, uid, cluster.Node)
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

// AskNode 检测用户是否在给定的节点上
func (p *Proxy) AskNode(ctx context.Context, uid int64, nid string) (string, bool, error) {
	if val, ok := p.nodeSource.Load(uid); ok {
		if val.(string) == nid {
			return nid, true, nil
		}
	}
	insID, err := p.node.opts.locator.Get(ctx, uid, cluster.Node)
	if err != nil {
		return "", false, err
	}
	if insID == "" {
		p.nodeSource.Delete(uid)
		return "", false, ErrNotFoundUserSource
	}
	p.nodeSource.Store(uid, insID)

	return insID, insID == nid, nil
}

// FetchGateList 拉取网关列表
func (p *Proxy) FetchGateList(ctx context.Context, states ...int32) ([]*registry.ServiceInstance, error) {
	return p.FetchServiceList(ctx, cluster.Gate, states...)
}

// FetchNodeList 拉取节点列表
func (p *Proxy) FetchNodeList(ctx context.Context, states ...int32) ([]*registry.ServiceInstance, error) {
	return p.FetchServiceList(ctx, cluster.Node, states...)
}

// 获得路由相关的node列表
func (p *Proxy) FetchNodeIdListByRoute(ctx context.Context, routeIds ...int32) ([]string, error) {
	nodes, err := p.FetchServiceList(ctx, cluster.Node, cluster.Work)
	if err != nil {
		return nil, nil
	}
	var nodeIds []string
	for _, nodeInfo := range nodes {
		for _, route := range nodeInfo.Routes {
			if slice.InSliceInt32(route.ID, routeIds) {
				nodeIds = append(nodeIds, nodeInfo.ID)
				break
			}
		}
	}
	return nodeIds, nil
}

// FetchServiceList 拉取服务列表
func (p *Proxy) FetchServiceList(ctx context.Context, kind string, states ...int32) ([]*registry.ServiceInstance, error) {
	services, err := p.node.opts.registry.Services(ctx, kind)
	if err != nil {
		return nil, err
	}

	if len(states) == 0 {
		return services, nil
	}

	mp := make(map[int32]struct{}, len(states))
	for _, state := range states {
		mp[state] = struct{}{}
	}

	list := make([]*registry.ServiceInstance, 0, len(services))
	for i := range services {
		if _, ok := mp[services[i].State]; ok {
			list = append(list, services[i])
		}
	}

	return list, nil
}

// GetIP 获取客户端IP
func (p *Proxy) GetIP(ctx context.Context, args *GetIPArgs) (string, error) {
	// GetIP 获取客户端IP
	switch args.Kind {
	case session.Conn:
		return p.directGetIP(ctx, args.GID, args.Kind, args.Target)
	case session.User:
		if args.GID == "" {
			return p.indirectGetIP(ctx, args.Target)
		} else {
			return p.directGetIP(ctx, args.GID, args.Kind, args.Target)
		}
	default:
		return "", ErrInvalidSessionKind
	}
}

// 直接获取IP
func (p *Proxy) directGetIP(ctx context.Context, gid string, kind session.Kind, target int64) (string, error) {
	client, err := p.getGateClientByGID(gid)
	if err != nil {
		return "", err
	}
	ip, _, err := client.GetIP(ctx, kind, target)
	return ip, err
}

// 间接获取IP
func (p *Proxy) indirectGetIP(ctx context.Context, uid int64) (string, error) {
	v, err := p.doGateRPC(ctx, uid, func(client *GateGrpcClient) (bool, interface{}, error) {
		ip, miss, err := client.GetIP(ctx, session.User, uid)
		return miss, ip, err
	})
	if err != nil {
		return "", err
	}

	return v.(string), nil
}

// Push 推送消息
func (p *Proxy) Push(ctx context.Context, args *PushArgs) error {
	switch args.Kind {
	case session.Conn:
		return p.directPush(ctx, args)
	case session.User:
		if args.GID == "" {
			return p.indirectPush(ctx, args)
		} else {
			return p.directPush(ctx, args)
		}
	default:
		return ErrInvalidSessionKind
	}
}

// 直接推送
func (p *Proxy) directPush(ctx context.Context, args *PushArgs) error {
	buffer, err := p.toBuffer(args.Message.Data)
	if err != nil {
		return err
	}
	client, err := p.getGateClientByGID(args.GID)
	if err != nil {
		return err
	}
	_, err = client.Push(ctx, args.Kind, args.Target, &packet.Message{
		Seq:    args.Message.Seq,
		Route:  args.Message.Route,
		Buffer: buffer,
	})
	return err
}

// 消息转buffer
func (l *Proxy) toBuffer(message interface{}) ([]byte, error) {
	if message == nil {
		return nil, nil
	}
	if v, ok := message.([]byte); ok {
		return v, nil
	}
	data, err := l.node.opts.codec.Marshal(message)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// 间接推送
func (l *Proxy) indirectPush(ctx context.Context, args *PushArgs) error {

	buffer, err := l.toBuffer(args.Message.Data)
	if err != nil {
		return err
	}
	_, err = l.doGateRPC(ctx, args.Target, func(client *GateGrpcClient) (bool, interface{}, error) {
		miss, pErr := client.Push(ctx, session.User, args.Target, &packet.Message{
			Seq:    args.Message.Seq,
			Route:  args.Message.Route,
			Buffer: buffer,
		})
		return miss, nil, pErr
	})
	return err
}

// Multicast 推送组播消息
func (p *Proxy) Multicast(ctx context.Context, args *MulticastArgs) (int64, error) {
	switch args.Kind {
	case session.Conn:
		return p.directMulticast(ctx, args)
	case session.User:
		if args.GID == "" {
			return p.indirectMulticast(ctx, args)
		} else {
			return p.directMulticast(ctx, args)
		}
	default:
		return 0, ErrInvalidSessionKind
	}
}

// 直接推送组播消息，只能推送到同一个网关服务器上
func (p *Proxy) directMulticast(ctx context.Context, args *MulticastArgs) (int64, error) {
	if len(args.Targets) == 0 {
		return 0, ErrReceiveTargetEmpty
	}
	buffer, err := p.toBuffer(args.Message.Data)
	if err != nil {
		return 0, err
	}
	client, err := p.getGateClientByGID(args.GID)
	if err != nil {
		return 0, err
	}
	return client.Multicast(ctx, args.Kind, args.Targets, &packet.Message{
		Seq:    args.Message.Seq,
		Route:  args.Message.Route,
		Buffer: buffer,
	})
}

// 间接推送组播消息
func (p *Proxy) indirectMulticast(ctx context.Context, args *MulticastArgs) (int64, error) {
	buffer, err := p.toBuffer(args.Message.Data)
	if err != nil {
		return 0, err
	}

	total := int64(0)
	eg, ctx := errgroup.WithContext(ctx)
	for _, target := range args.Targets {
		func(target int64) {
			eg.Go(func() error {
				_, dErr := p.doGateRPC(ctx, target, func(client *GateGrpcClient) (bool, interface{}, error) {
					miss, pErr := client.Push(ctx, session.User, target, &packet.Message{
						Seq:    args.Message.Seq,
						Route:  args.Message.Route,
						Buffer: buffer,
					})
					return miss, nil, pErr
				})
				if dErr != nil {
					return dErr
				}

				atomic.AddInt64(&total, 1)
				return nil
			})
		}(target)
	}

	err = eg.Wait()

	if total > 0 {
		return total, nil
	}

	return 0, err
}

// Broadcast 推送广播消息
func (p *Proxy) Broadcast(ctx context.Context, args *BroadcastArgs) (int64, error) {
	buffer, err := p.toBuffer(args.Message.Data)
	if err != nil {
		return 0, err
	}
	total := int64(0)
	eg, ctx := errgroup.WithContext(ctx)
	p.gateDispatcher.IterateEndpoint(func(_ string, ep *endpoint.Endpoint) bool {
		eg.Go(func() error {
			client, nErr := p.node.opts.transporter.NewGateClient(ep)
			if nErr != nil {
				return nErr
			}
			n, bErr := client.Broadcast(ctx, args.Kind, &packet.Message{
				Seq:    args.Message.Seq,
				Route:  args.Message.Route,
				Buffer: buffer,
			})
			if bErr != nil {
				return bErr
			}
			atomic.AddInt64(&total, n)
			return nil
		})
		return true
	})
	err = eg.Wait()

	if total > 0 {
		return total, nil
	}
	return total, err
}

// Deliver投递消息给节点处理
func (p *Proxy) Deliver(ctx context.Context, args *DeliverArgs) error {
	if args.NID != p.GetNodeID() {
		return p.Deliver2(ctx, &DeliverArgs2{
			NID: args.NID,
			UID: args.UID,
			Message: &Message{
				Seq:   args.Message.Seq,
				Route: args.Message.Route,
				Data:  args.Message.Data,
			},
		})
	} else {
		p.node.router.deliver("", args.NID, 0, args.UID, args.Message.Seq, args.Message.Route, args.Message.Data)
	}
	return nil
}

// Deliver 投递消息给节点处理
func (p *Proxy) Deliver2(ctx context.Context, args *DeliverArgs2) error {
	arguments := &cluster.DeliverArgs{
		GID: "",
		NID: p.node.opts.id,
		CID: args.CID,
		UID: args.UID,
	}

	switch msg := args.Message.(type) {
	case *packet.Message:
		arguments.Message = &packet.Message{
			Seq:    msg.Seq,
			Route:  msg.Route,
			Buffer: msg.Buffer,
		}
	case *Message:
		buffer, err := p.toBuffer(msg.Data)
		if err != nil {
			return err
		}
		arguments.Message = &packet.Message{
			Seq:    msg.Seq,
			Route:  msg.Route,
			Buffer: buffer,
		}
	default:
		return ErrInvalidMessage
	}

	if args.NID != "" {
		client, err := p.getNodeClientByNID(args.NID)
		if err != nil {
			return err
		}
		_, err = client.Deliver(ctx, arguments)
		return err
	} else {
		_, err := p.doNodeRPC(ctx, arguments.Message.Route, args.UID, func(ctx context.Context, client *NodeGrpcClient) (bool, interface{}, error) {
			miss, err := client.Deliver(ctx, arguments)
			return miss, nil, err
		})
		return err
	}
}

// 执行节点RPC调用
func (p *Proxy) doNodeRPC(ctx context.Context, routeID int32, uid int64, fn func(ctx context.Context, client *NodeGrpcClient) (bool, interface{}, error)) (interface{}, error) {
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

		client, err = p.node.opts.transporter.NewNodeClient(ep)
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

// Response 响应消息
func (p *Proxy) Response(ctx context.Context, req *Request, data interface{}) error {
	switch {
	// 网关过来的消息
	case req.GID != "":
		return p.Push(ctx, &PushArgs{
			GID:    req.GID,
			Kind:   session.Conn,
			Target: req.CID,
			Message: &Message{
				Seq:   req.Message.Seq,
				Route: req.Message.Route,
				Data:  data,
			},
		})
		// 节点 给节点投递的消息
	case req.NID != "":
		return p.Deliver(ctx, &DeliverArgs{
			NID: req.NID,
			UID: req.UID,
			Message: &Message{
				Seq:   req.Message.Seq,
				Route: req.Message.Route,
				Data:  data,
			},
		})
	}

	return nil
}

// Disconnect 断开连接
func (p *Proxy) Disconnect(ctx context.Context, args *DisconnectArgs) error {
	switch args.Kind {
	case session.Conn:
		return p.directDisconnect(ctx, args.GID, args.Kind, args.Target, args.IsForce)
	case session.User:
		if args.GID == "" {
			return p.indirectDisconnect(ctx, args.Target, args.IsForce)
		} else {
			return p.directDisconnect(ctx, args.GID, args.Kind, args.Target, args.IsForce)
		}
	default:
		return ErrInvalidSessionKind
	}
}

// 直接断开连接
func (p *Proxy) directDisconnect(ctx context.Context, gid string, kind session.Kind, target int64, isForce bool) error {
	client, err := p.getGateClientByGID(gid)
	if err != nil {
		return err
	}
	_, err = client.Disconnect(ctx, kind, target, isForce)
	return err
}

// 间接断开连接
func (p *Proxy) indirectDisconnect(ctx context.Context, uid int64, isForce bool) error {
	_, err := p.doGateRPC(ctx, uid, func(client *GateGrpcClient) (bool, interface{}, error) {
		miss, err := client.Disconnect(ctx, session.User, uid, isForce)
		return miss, nil, err
	})
	return err
}

// Invoke 调用函数 （可以理解为单线程操作)
//func (p *Proxy) Invoke(fn func()) {
//	p.node.fnChan <- fn
//}

// 启动监听
func (p *Proxy) watch(ctx context.Context) {
	p.watchUserLocate(ctx, cluster.Gate, cluster.Node)
	p.WatchServiceInstance(ctx, cluster.Gate, cluster.Node)
}

// WatchServiceInstance 监听服务实例
func (p *Proxy) WatchServiceInstance(ctx context.Context, kinds ...string) {
	for _, kind := range kinds {
		p.watchServiceInstance(ctx, kind)
	}
}

// 监听服务实例
func (p *Proxy) watchServiceInstance(ctx context.Context, kind string) {
	rctx, rcancel := context.WithTimeout(ctx, 10*time.Second)
	watcher, err := p.node.opts.registry.Watch(rctx, kind)
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
			if kind == cluster.Node {
				p.nodeDispatcher.ReplaceServices(services...)
			} else {
				p.gateDispatcher.ReplaceServices(services...)
			}
		}
	}()
}

// watchUserLocate 监听用户定位
func (p *Proxy) watchUserLocate(ctx context.Context, kinds ...string) {
	rctx, rcancel := context.WithTimeout(ctx, 10*time.Second)
	watcher, err := p.node.opts.locator.Watch(rctx, kinds...)
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
				case cluster.Gate:
					source = &p.gateSource
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
