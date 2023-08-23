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
	"github.com/cr-mao/loric/log"
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
	// 监听重新连接
	r.proxy.Events().AddEventHandler(cluster.Reconnect, r.reconnect)
	// 监听连接断开
	r.proxy.Events().AddEventHandler(cluster.Disconnect, r.disconnect)
	// 创建路由
	r.proxy.Router().AddRouteHandler(Login, false, r.Login)
}

func (r *Router) Login(ctx *node.Context) {
	log.Infof("gid:%s,nid:%s,cid:%d,uid:%d", ctx.Request.GID, ctx.Request.NID, ctx.Request.CID, ctx.Request.UID)
	log.Info(ctx.Request.Message)

}

// 重新连接
func (r *Router) reconnect(evt *node.Event) {
	log.Warnf("connection is reopened, gid: %v, cid: %d, uid: %d", evt.GID, evt.CID, evt.UID)
}

// 连接断开
func (r *Router) disconnect(evt *node.Event) {
	log.Warnf("connection is closed, gid: %v, cid: %d, uid: %d", evt.GID, evt.CID, evt.UID)
}
