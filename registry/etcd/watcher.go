package etcd

import (
	"context"
	"strings"
	"sync"
	"sync/atomic"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/cr-mao/loric/registry"
)

type watcherMgr struct {
	err              error
	ctx              context.Context
	cancel           context.CancelFunc
	registry         *Registry
	serviceName      string
	serviceInstances sync.Map
	watcher          clientv3.Watcher
	chWatch          clientv3.WatchChan

	idx      int64
	rw       sync.RWMutex
	watchers map[int64]*watcher
}

type watcher struct {
	idx        int64
	state      int32
	watcherMgr *watcherMgr
	ctx        context.Context
	cancel     context.CancelFunc
	chWatch    chan []*registry.ServiceInstance
}

func newWatcher(wm *watcherMgr, idx int64) *watcher {
	w := &watcher{}
	w.ctx, w.cancel = context.WithCancel(wm.ctx)
	w.idx = idx
	w.watcherMgr = wm
	w.chWatch = make(chan []*registry.ServiceInstance, 16)

	return w
}

func (w *watcher) notify(services []*registry.ServiceInstance) {
	if atomic.LoadInt32(&w.state) == 0 {
		return
	}

	w.chWatch <- services
}

// Next 返回服务实例列表
func (w *watcher) Next() ([]*registry.ServiceInstance, error) {
	if atomic.LoadInt32(&w.state) == 0 {
		atomic.StoreInt32(&w.state, 1)
		return w.watcherMgr.services(), nil
	}

	select {
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	case services, ok := <-w.chWatch:
		if !ok {
			if err := w.ctx.Err(); err != nil {
				return nil, err
			}
		}

		return services, nil
	}
}

// Stop 停止监听
func (w *watcher) Stop() error {
	w.cancel()
	close(w.chWatch)
	return w.watcherMgr.recycle(w.idx)
}

func newWatcherMgr(registry *Registry, ctx context.Context, serviceName string) (*watcherMgr, error) {
	services, err := registry.services(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	w := &watcherMgr{}
	w.ctx, w.cancel = context.WithCancel(registry.ctx)
	w.registry = registry
	w.serviceName = serviceName
	w.watcher = clientv3.NewWatcher(registry.opts.client)
	w.chWatch = w.watcher.Watch(w.ctx, buildPrefixKey(registry.opts.namespace, w.serviceName), clientv3.WithPrefix())
	w.watchers = make(map[int64]*watcher)

	for _, service := range services {
		w.serviceInstances.Store(service.ID, service)
	}

	go func() {
		for {
			select {
			case <-w.ctx.Done():
				return
			case res := <-w.chWatch:
				for _, ev := range res.Events {
					switch ev.Type {
					case mvccpb.PUT:
						if service, err := unmarshal(ev.Kv.Value); err == nil {
							w.serviceInstances.Store(service.ID, service)
						}
					case mvccpb.DELETE:
						if parts := strings.Split(string(ev.Kv.Key), "/"); len(parts) == 4 {
							w.serviceInstances.Delete(parts[3])
						}
					}
				}

				w.broadcast()
			}
		}
	}()

	return w, nil
}

func (wm *watcherMgr) fork() registry.Watcher {
	wm.rw.Lock()
	defer wm.rw.Unlock()

	w := newWatcher(wm, atomic.AddInt64(&wm.idx, 1))
	wm.watchers[w.idx] = w

	return w
}

func (wm *watcherMgr) recycle(idx int64) error {
	wm.rw.Lock()
	defer wm.rw.Unlock()

	delete(wm.watchers, idx)

	if len(wm.watchers) == 0 {
		wm.cancel()
		wm.registry.watchers.Delete(wm.serviceName)
		return wm.watcher.Close()
	}

	return nil
}

func (wm *watcherMgr) broadcast() {
	wm.rw.RLock()
	defer wm.rw.RUnlock()

	services := wm.services()
	for _, w := range wm.watchers {
		w.notify(services)
	}
}

func (wm *watcherMgr) services() (services []*registry.ServiceInstance) {
	wm.serviceInstances.Range(func(key, value interface{}) bool {
		services = append(services, value.(*registry.ServiceInstance))
		return true
	})
	return
}
