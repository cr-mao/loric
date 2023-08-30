/**
User: cr-mao
Date: 2023/8/24 13:52
Email: crmao@qq.com
Desc: event.go
*/
package controller

import (
	"context"
	"github.com/cr-mao/loric/cluster/node"
	"github.com/cr-mao/loric/log"
	"time"
)

type EventController struct{}

// 重新连接
func (r *EventController) Reconnect(evt *node.Event) {
	log.Infof("connection is reopened, gid: %v, cid: %d, uid: %d", evt.GID, evt.CID, evt.UID)
}

// 连接断开
func (r *EventController) Disconnect(evt *node.Event) {
	log.Infof("connection is closed, gid: %v, cid: %d, uid: %d", evt.GID, evt.CID, evt.UID)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	err := evt.Proxy.UnbindNode(ctx, evt.UID, evt.Proxy.GetNodeID())
	if err != nil {
		log.Errorf("event disconnect err %v", err)
	}
}
