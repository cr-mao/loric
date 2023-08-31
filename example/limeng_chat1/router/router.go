/**
User: cr-mao
Date: 2023/8/23 15:25
Email: crmao@qq.com
Desc: router.go
*/
package router

import (
	"github.com/cr-mao/loric/cluster"
	"github.com/cr-mao/loric/cluster/node"
	"github.com/cr-mao/loric/example/internal/pb"
	"github.com/cr-mao/loric/example/limeng_chat/controller"
)

type Router struct {
	proxy *node.Proxy
}

func NewRouter(proxy *node.Proxy) *Router {
	return &Router{
		proxy: proxy,
	}
}

func (r *Router) Init() {
	var eventController = &controller.EventController{}
	var msgController = &controller.MsgController{}
	// 监听重新连接
	r.proxy.Events().AddEventHandler(cluster.Reconnect, eventController.Reconnect)
	// 监听连接断开
	r.proxy.Events().AddEventHandler(cluster.Disconnect, eventController.Disconnect)
	// 创建路由
	r.proxy.Router().AddRouteHandler(int32(pb.Route_LianmentChatEnter), true, msgController.Enter)
	r.proxy.Router().AddRouteHandler(int32(pb.Route_LianmengChat), true, msgController.MsgHandle)
}
