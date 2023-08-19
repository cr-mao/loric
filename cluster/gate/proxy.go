package gate

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cr-mao/loric/cluster"
	"github.com/cr-mao/loric/errors"
	"github.com/cr-mao/loric/locate"
	"github.com/cr-mao/loric/log"
	"github.com/cr-mao/loric/packet"
)

var (
	ErrInvalidArgument = errors.New("ErrInvalidArgument")
)

type proxy struct {
	gate       *Gate    // 网关
	nodeSource sync.Map //  用户在哪台 node

}

func newProxy(gate *Gate) *proxy {
	return &proxy{gate: gate}
}

// 绑定用户与网关间的关系
func (p *proxy) bindGate(ctx context.Context, cid, uid int64) error {
	err := p.gate.opts.locator.Set(ctx, uid, cluster.Gate, p.gate.opts.id)
	if err != nil {
		return err
	}
	// 绑定 网关算重连， 登录的时候 ，node过来触发的
	p.trigger(ctx, cluster.Reconnect, cid, uid)
	return nil
}

// 解绑用户与网关间的关系
func (p *proxy) unbindGate(ctx context.Context, cid, uid int64) error {
	err := p.gate.opts.locator.Rem(ctx, uid, cluster.Gate, p.gate.opts.id)
	if err != nil {
		log.Errorf("user unbind failed, gid: %d, cid: %d, uid: %d, err: %v", p.gate.opts.id, cid, uid, err)
	}
	return err
}

// 触发事件
func (p *proxy) trigger(ctx context.Context, event cluster.Event, cid, uid int64) {
	fmt.Println("trigger", event, time.Now().UnixMilli())

	//if err := p.link.Trigger(ctx, &link.TriggerArgs{
	//	Event: event,
	//	CID:   cid,
	//	UID:   uid,
	//}); err != nil {
	//	log.Warnf("trigger event failed, gid: %s, cid: %d, uid: %d, event: %v, err: %v", p.gate.opts.id, cid, uid, event, err)
	//}

}

// 投递消息
func (p *proxy) deliver(ctx context.Context, cid, uid int64, data []byte) {
	message, err := packet.Unpack(data)
	if err != nil {
		log.Errorf("unpack data to struct failed: %v", err)
		return
	}

	fmt.Println(message)

	//if err = p.link.Deliver(ctx, &link.DeliverArgs{
	//	CID:     cid,
	//	UID:     uid,
	//	Message: message,
	//}); err != nil {
	//	log.Errorf("deliver message failed: %v", err)
	//}
}

// 启动监听
func (p *proxy) watch(ctx context.Context) {
	p.watchUserLocate(ctx)
	p.watchServiceInstance(ctx)
}

// 监听服务实例
func (p *proxy) watchServiceInstance(ctx context.Context) {
	rctx, rcancel := context.WithTimeout(ctx, 10*time.Second)
	watcher, err := p.gate.opts.registry.Watch(rctx, string(cluster.Node))
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
			//services, err := watcher.Next()
			//if err != nil {
			//	continue
			//}
			//p.nodeDispatcher.ReplaceServices(services...)
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
