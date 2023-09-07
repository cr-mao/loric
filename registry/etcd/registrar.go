package etcd

import (
	"context"
	"errors"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/cr-mao/loric/log"
	"github.com/cr-mao/loric/registry"
)

type heartbeat struct {
	leaseID clientv3.LeaseID
	key     string
	value   string
}

type registrar struct {
	registry    *Registry
	ctx         context.Context
	cancel      context.CancelFunc
	kv          clientv3.KV
	lease       clientv3.Lease
	chHeartbeat chan heartbeat
}

func newRegistrar(registry *Registry) *registrar {
	r := &registrar{}
	r.kv = clientv3.NewKV(registry.opts.client)
	r.lease = clientv3.NewLease(registry.opts.client)
	r.ctx, r.cancel = context.WithCancel(registry.ctx)
	r.registry = registry
	r.chHeartbeat = make(chan heartbeat)

	go func() {
		var (
			ctx    context.Context
			cancel context.CancelFunc
		)

		for {
			select {
			case heartbeats, ok := <-r.chHeartbeat:
				if cancel != nil {
					cancel()
				}

				if !ok {
					return
				}

				ctx, cancel = context.WithCancel(r.ctx)
				go r.heartbeat(ctx, heartbeats.leaseID, heartbeats.key, heartbeats.value)
			case <-r.ctx.Done():
				if cancel != nil {
					cancel()
				}
				return
			}
		}
	}()

	return r
}

// 注册服务
func (r *registrar) register(ctx context.Context, ins *registry.ServiceInstance) error {
	value, err := marshal(ins)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("/%s/%s/%s", r.registry.opts.namespace, ins.Name, ins.ID)

	leaseID, err := r.put(ctx, key, value)
	if err != nil {
		return err
	}

	r.chHeartbeat <- heartbeat{
		leaseID: leaseID,
		key:     key,
		value:   value,
	}

	return nil
}

// 解注册服务
func (r *registrar) deregister(ctx context.Context, ins *registry.ServiceInstance) (err error) {
	r.cancel()
	close(r.chHeartbeat)

	r.registry.registrars.Delete(ins.ID)

	key := fmt.Sprintf("/%s/%s/%s", r.registry.opts.namespace, ins.Name, ins.ID)
	_, err = r.kv.Delete(ctx, key)

	//nolint
	if r.lease != nil {
		_ = r.lease.Close()
	}

	return
}

// 写入KV
func (r *registrar) put(ctx context.Context, key, value string) (clientv3.LeaseID, error) {
	res, err := r.lease.Grant(ctx, int64(r.registry.opts.retryInterval.Seconds())+2)
	if err != nil {
		return 0, err
	}
	_, err = r.kv.Put(ctx, key, value, clientv3.WithLease(res.ID))
	if err != nil {
		return 0, err
	}

	return res.ID, nil
}

func (r *registrar) getByKey(ctx context.Context, key string) (string, error) {
	var getResp *clientv3.GetResponse
	getResp, err := r.kv.Get(ctx, key)
	if err != nil {
		return "", err
	}
	if len(getResp.Kvs) >= 1 {
		return string(getResp.Kvs[0].Value), nil
	}
	return "", errors.New("instance not found")
}

// 心跳
func (r *registrar) heartbeat(ctx context.Context, leaseID clientv3.LeaseID, key, value string) {
	chKA, err := r.lease.KeepAlive(ctx, leaseID)
	ok := err == nil

	for {
		if !ok {
			for i := 0; i < r.registry.opts.retryTimes; i++ {
				log.Warn("retry keepalive")
				if ctx.Err() != nil {
					return
				}
				pctx, pcancel := context.WithTimeout(ctx, r.registry.opts.timeout)
				// 重新取一遍，万一失败，则取第一次的value了。
				newValue, kerr := r.getByKey(ctx, key)
				if kerr != nil {
					newValue = value
				}
				// 管理服，把权重给改了。 如果这里keepalive 失败 ，那么会写入 老的值（ 按理可能性不大））
				leaseID, err = r.put(pctx, key, newValue)
				pcancel()
				if err != nil {
					time.Sleep(r.registry.opts.retryInterval)
					continue
				}

				chKA, err = r.lease.KeepAlive(ctx, leaseID)
				if err != nil {
					time.Sleep(r.registry.opts.retryInterval)
					continue
				}
				ok = true
				break
			}
			if !ok {
				return
			}
		}

		select {
		case _, ok = <-chKA:
			//log.Debug("keep alive ", ok, time.Now().Unix())
			if !ok {
				if ctx.Err() != nil {
					return
				}
				continue
			}
		case <-ctx.Done():
			return
		}
	}
}
