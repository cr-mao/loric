package registry

import (
	"context"
)

type Registry interface {
	// Register 注册服务实例
	Register(ctx context.Context, ins *ServiceInstance) error
	// Deregister 解注册服务实例
	Deregister(ctx context.Context, ins *ServiceInstance) error
	// Watch 监听相同服务名的服务实例变化
	Watch(ctx context.Context, serviceName string) (Watcher, error)
	// Services 获取服务实例列表
	Services(ctx context.Context, serviceName string) ([]*ServiceInstance, error)
}

type Discovery interface {
	// Watch 监听相同服务名的服务实例变化
	Watch(ctx context.Context, serviceName string) (Watcher, error)
	// Services 获取服务实例列表
	Services(ctx context.Context, serviceName string) ([]*ServiceInstance, error)
}

type Watcher interface {
	// Next 返回服务实例列表
	Next() ([]*ServiceInstance, error)
	// Stop 停止监听
	Stop() error
}

type ServiceInstance struct {
	// 服务实体ID，每个服务实体ID唯一
	ID string `json:"id"`
	// 服务实体名
	Name string `json:"name"`
	// 服务实体类型
	Kind string `json:"kind"`
	// 服务实体别名
	Alias string `json:"alias"`
	// 服务实例状态
	State int32 `json:"state"`
	// 服务事件集合
	Events []int32 `json:"events"`
	// 服务路由ID
	Routes []Route `json:"routes"`
	// 服务器实体暴露端口
	Endpoint string `json:"endpoint"`
	// 服务权重， 权重负载均衡，可以设置为0 ，则表示新的连接 不会请求过去，也可以用来做热更用。
	Weight int `json:"weight"`
}

type Route struct {
	// 路由ID
	ID int32 `json:"id"`
	// 是否有状态
	Stateful bool `json:"stateful"`
}
