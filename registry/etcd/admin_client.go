/**
User: cr-mao
Date: 2023/11/30 14:03
Email: crmao@qq.com
Desc: 给管理后台用的。框架这边不用
*/
package etcd

import (
	"context"

	"github.com/cr-mao/loric/encoding/json"
	"github.com/cr-mao/loric/errors"
	"github.com/cr-mao/loric/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// 	key := fmt.Sprintf("/%s/%s/%s", namespace,  gate|node  , 服务唯一id)
func (c *Registry) GetByKey(ctx context.Context, key string) (*registry.ServiceInstance, error) {
	kv := clientv3.NewKV(c.opts.client)
	ctx, cancel := context.WithTimeout(ctx, c.opts.timeout)
	defer cancel()
	var getResp *clientv3.GetResponse
	getResp, err := kv.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	var instance *registry.ServiceInstance
	for _, v := range getResp.Kvs {
		instance, err = unmarshal(v.Value)
		if err != nil {
			return nil, err
		}
		return instance, nil
	}
	return nil, errors.New("instance not found")
}

// 设置、覆盖值，如设置权重
func (c *Registry) Pub(ctx context.Context, key string, value string) error {
	kv := clientv3.NewKV(c.opts.client)
	_, err := kv.Put(ctx, key, value)
	return err
}

// namespace 命名空间，  serviceName : gate|node
func (c *Registry) GetServerInstances(ctx context.Context, namespace, serviceName string) (res []*registry.ServiceInstance, err error) {
	key := buildPrefixKey(namespace, serviceName)
	kv := clientv3.NewKV(c.opts.client)
	ctx, cancel := context.WithTimeout(ctx, c.opts.timeout)
	defer cancel()
	var getResp *clientv3.GetResponse
	getResp, err = kv.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	var instance *registry.ServiceInstance
	for _, v := range getResp.Kvs {
		instance, err = unmarshal(v.Value)
		if err != nil {
			return
		}
		res = append(res, instance)
	}
	return
}

func Marshal(ins *registry.ServiceInstance) (string, error) {
	buf, err := json.Marshal(ins)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
