/**
User: cr-mao
Date: 2023/9/5 18:28
Email: crmao@qq.com
Desc: etcd.go
*/
package etcdclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/cr-mao/loric/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdClient struct {
	// 内建客户端配置，默认为[]string{"localhost:2379"}
	addrs []string
	// 客户端拨号超时时间
	// 内建客户端配置，默认为5秒
	dialTimeout time.Duration
	// 外部客户端
	// 外部客户端配置，存在外部客户端时，优先使用外部客户端，默认为nil
	client *clientv3.Client
	// 上下文
	// 默认为3秒
	timeout time.Duration

	err error
}

var GlobalEtcdClient *EtcdClient

// 内网不需要账号，密码
func NewEtcdClient(addrs []string, dialTimeout time.Duration) *EtcdClient {
	r := &EtcdClient{
		timeout: 3 * time.Second,
	}
	r.client, r.err = clientv3.New(clientv3.Config{
		Endpoints:   addrs,
		DialTimeout: dialTimeout,
	})

	GlobalEtcdClient = r
	return r
}

func buildPrefixKey(namespace, serviceName string) string {
	return fmt.Sprintf("/%s/%s", namespace, serviceName)
}

func (c *EtcdClient) GetByKey(ctx context.Context, key string) (*registry.ServiceInstance, error) {
	kv := clientv3.NewKV(c.client)
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
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

func (c *EtcdClient) Pub(ctx context.Context, key string, value string) error {
	kv := clientv3.NewKV(c.client)
	_, err := kv.Put(ctx, key, value)
	return err
}

func (c *EtcdClient) GetServerInstances(ctx context.Context, namespace, serviceName string) (res []*registry.ServiceInstance, err error) {
	key := buildPrefixKey(namespace, serviceName)
	kv := clientv3.NewKV(c.client)
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
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

func unmarshal(data []byte) (*registry.ServiceInstance, error) {
	ins := &registry.ServiceInstance{}
	if err := json.Unmarshal(data, ins); err != nil {
		return nil, err
	}
	return ins, nil
}
