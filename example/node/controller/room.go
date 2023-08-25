/**
User: cr-mao
Date: 2023/8/24 19:19
Email: crmao@qq.com
Desc: room.go
*/
package controller

import (
	"github.com/cr-mao/loric/cluster/node"
	"github.com/cr-mao/loric/example/internal/pb"
	"github.com/cr-mao/loric/log"
)

type RoomController struct{}

func (r *RoomController) CreateRoom(ctx *node.Context) {
	log.Infof("CreateRoom gid:%s,nid:%s,cid:%d,uid:%d", ctx.Request.GID, ctx.Request.NID, ctx.Request.CID, ctx.Request.UID)
	req := &pb.CreateRoomReq{}
	res := &pb.CreateRoomRes{}
	defer func() {
		if err := ctx.Response(res); err != nil {
			log.Errorf("response login message failed, err: %v", err)
		}
	}()
	if err := ctx.Request.Parse(req); err != nil {
		log.Errorf("invalid login message, err: %v", err)
		res.ID = 0
		return
	}
	if err := ctx.BindNode(); err != nil {
		log.Errorf("bind node failed, err: %v", err)
		res.Code = pb.CreateRoomCode_Failed
		return
	}
	res.Code = pb.CreateRoomCode_Ok
	// 强行返回房间1。
	res.ID = 1
}
