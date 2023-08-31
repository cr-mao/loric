/**
User: cr-mao
Date: 2023/8/24 13:54
Email: crmao@qq.com
Desc: login.go
*/
package controller

import (
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/cr-mao/loric/cluster/node"
	"github.com/cr-mao/loric/example/auth/global"
	"github.com/cr-mao/loric/example/internal/pb"
	"github.com/cr-mao/loric/log"
	"github.com/cr-mao/loric/sugar"
)

type LoginController struct{}

func (r *LoginController) Login(ctx *node.Context) {
	log.Infof("LoginHandle gid:%s,nid:%s,cid:%d,uid:%d", ctx.Request.GID, ctx.Request.NID, ctx.Request.CID, ctx.Request.UID)
	req := &pb.LoginReq{}
	res := &pb.LoginRes{}
	defer func() {
		if err := ctx.Response(res); err != nil {
			log.Errorf("response login message failed, err: %v", err)
		}
	}()

	res.Code = pb.LoginCode_Failed
	if err := ctx.Request.Parse(req); err != nil {
		log.Errorf("invalid login message, err: %v", err)
		return
	}
	userInfo, ok := global.UserData.UserMap.Load(req.Token)
	var userRow *global.User
	if !ok {
		newUserId := atomic.AddInt64(&global.UserId, 1)
		userName, _ := sugar.UUID()
		userName = userName[0:6]
		userRow = &global.User{
			UserId:   newUserId,
			UserName: userName,
			LiMengId: newUserId % 2, // 0,1
		}
		global.UserData.UserMap.Store(req.Token, userRow)
	} else {
		userRow = userInfo.(*global.User)
	}
	res.Code = pb.LoginCode_Ok
	if err := ctx.BindGate(userRow.UserId); err != nil {
		log.Errorf("bind gate failed: %v", err)
		return
	}
	nodes, err := ctx.Proxy.FetchNodeIdListByRoute(ctx.Context(), int32(pb.Route_LianmengChat))
	if err != nil {
		log.Errorf("login FetchNodeIdListByRoute err :%+v", err)
		return
	}
	if len(nodes) == 0 {
		log.Warnf("lianmeng chat node is not exist")
		return
	}

	// 暂时不知道咋解决先
	// lianmengID => nodeId
	// 绑定node ， 联盟操作的node ....
	lianmengIdStr := fmt.Sprintf("%d", userRow.LiMengId)
	var bandNodeId string
	for _, nodeId := range nodes {
		if strings.Contains(nodeId, lianmengIdStr) {
			bandNodeId = nodeId
			break
		}
	}
	//避免这种事情发生，nodes是无规律顺序的
	if bandNodeId == "" {
		bandNodeId = nodes[0]
	}
	fmt.Println("bind_node_id:", bandNodeId)
	fmt.Println("user_id:", userRow.UserId)
	err = ctx.Proxy.BindNode(ctx.Context(), userRow.UserId, bandNodeId)
	if err != nil {
		log.Errorf("login bind node err:%+v", err)
		return
	}
	res.Code = pb.LoginCode_Ok
}
