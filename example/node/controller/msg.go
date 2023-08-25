/**
User: cr-mao
Date: 2023/8/25 14:04
Email: crmao@qq.com
Desc: msg.go
*/
package controller

import (
	"github.com/cr-mao/loric/cluster/node"
	"github.com/cr-mao/loric/example/internal/pb"
	"github.com/cr-mao/loric/log"
)

type MsgController struct{}

func (c *MsgController) MsgHandle(ctx *node.Context) {
	//log.Infof("msgHandle gid:%s,nid:%s,cid:%d,uid:%d", ctx.Request.GID, ctx.Request.NID, ctx.Request.CID, ctx.Request.UID)
	req := &pb.SendMsgReq{}
	res := &pb.SendMsgRes{}
	defer func() {
		if err := ctx.Response(res); err != nil {
			log.Errorf("response login message failed, err: %v", err)
		}
	}()
	if err := ctx.Request.Parse(req); err != nil {
		log.Errorf("invalid login message, err: %v", err)
		res.Code = 0
		return
	}
	res.Msg = req.Msg + " from node:" + ctx.GetNodeId()
	if ctx.Request.UID == 1 {
		res.UserName = "cr-mao"
	}
	res.Code = 1
}
