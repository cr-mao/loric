/**
User: cr-mao
Date: 2023/8/24 13:54
Email: crmao@qq.com
Desc: login.go
*/
package controller

import (
	"github.com/cr-mao/loric/cluster/node"
	"github.com/cr-mao/loric/example/internal/pb"
	"github.com/cr-mao/loric/log"
)

type LoginController struct{}

func (r *LoginController) Login(ctx *node.Context) {
	//log.Infof("LoginHandle gid:%s,nid:%s,cid:%d,uid:%d", ctx.Request.GID, ctx.Request.NID, ctx.Request.CID, ctx.Request.UID)
	req := &pb.LoginReq{}
	res := &pb.LoginRes{}
	defer func() {
		if err := ctx.Response(res); err != nil {
			log.Errorf("response login message failed, err: %v", err)
		}
	}()
	if err := ctx.Request.Parse(req); err != nil {
		log.Errorf("invalid login message, err: %v", err)
		res.Code = pb.LoginCode_Failed
		return
	}
	var uid int64
	if req.Token == "cr-mao" {
		res.Code = pb.LoginCode_Ok
		uid = ctx.Request.CID //用户id 就是连接id 先
	}
	if err := ctx.BindGate(uid); err != nil {
		log.Errorf("bind gate failed: %v", err)
		return
	}

}
