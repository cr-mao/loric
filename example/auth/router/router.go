/**
User: cr-mao
Date: 2023/8/23 15:25
Email: crmao@qq.com
Desc: router.go
*/
package router

import (
	"github.com/cr-mao/loric/cluster/node"
	"github.com/cr-mao/loric/example/auth/controller"
	"github.com/cr-mao/loric/example/internal/pb"
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
	var loginController = &controller.LoginController{}
	// 创建路由
	r.proxy.Router().AddRouteHandler(int32(pb.Route_Login), false, loginController.Login)
}
