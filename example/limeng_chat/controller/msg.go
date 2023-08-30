/**
User: cr-mao
Date: 2023/8/25 14:04
Email: crmao@qq.com
Desc: msg.go
*/
package controller

import (
	"fmt"
	"github.com/cr-mao/loric/cluster/node"
	"github.com/cr-mao/loric/example/internal/pb"
	"github.com/cr-mao/loric/example/limeng_chat/manager"
	"github.com/cr-mao/loric/log"
	"github.com/cr-mao/loric/session"
)

type MsgController struct{}

func (c *MsgController) Enter(ctx *node.Context) {
	res := &pb.LianmengEnterResponse{}
	res.Code = pb.LianmengEnterCode_Failed
	defer func() {
		if err := ctx.Response(res); err != nil {
			log.Errorf("response login message failed, err: %v", err)
		}
	}()
	var user = &manager.User{
		UserId:     ctx.Request.UID,
		UserName:   fmt.Sprintf("user_%d", ctx.Request.UID),
		LianmengId: getLianmengId(ctx.Request.UID),
	}
	manager.GetManager().AddUser(user)
	res.Code = pb.LianmengEnterCode_Ok
}

func (c *MsgController) MsgHandle(ctx *node.Context) {
	userId := ctx.Request.UID
	log.Infof("msgHandle gid:%s,nid:%s,cid:%d,uid:%d", ctx.Request.GID, ctx.Request.NID, ctx.Request.CID, userId)
	req := &pb.LianmengChatMsgReq{}
	res := &pb.LianmengChatSendMsgRes{}
	res.Code = pb.LianmengChatCode_Failed
	defer func() {
		if err := ctx.Response(res); err != nil {
			log.Errorf("response login message failed, err: %v", err)
		}
	}()
	if err := ctx.Request.Parse(req); err != nil {
		log.Errorf("invalid  message, err: %v", err)
		return
	}
	user, err := manager.GetManager().GetUserByUserIdAndLianMengId(userId, getLianmengId(userId))
	if err != nil {
		log.Errorf("GetUserByUserIdAndLianMengId  , err: %v", err)
		return
	}
	res.Msg = "node_id:" + ctx.GetNodeId() + "get msg:" + req.Msg + " user_lianmeng_id:" + fmt.Sprintf("%d", user.LianmengId)
	res.UserName = user.UserName
	res.Code = pb.LianmengChatCode_Ok
	// 广播消息
	otherLianmengUserIds := manager.GetManager().GetUserIdListByLianmengId(getLianmengId(userId), userId)
	if len(otherLianmengUserIds) == 0 {
		return
	}
	fmt.Println(otherLianmengUserIds)
	mcnt, err := ctx.Proxy.Multicast(ctx.Context(), &node.MulticastArgs{
		GID:     ctx.Request.GID,
		Kind:    session.User,
		Targets: otherLianmengUserIds,
		Message: &node.Message{
			Seq:   ctx.Request.Message.Seq,
			Route: ctx.Request.Message.Route,
			Data:  res,
		},
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(mcnt)

}
