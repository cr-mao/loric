/**
User: cr-mao
Date: 2023/8/21 20:47
Email: crmao@qq.com
Desc: request.go
*/
package node

import (
	"github.com/jinzhu/copier"
)

// Request 请求数据
type Request struct {
	node    *Node
	GID     string   // 来源网关ID
	NID     string   // 来源节点ID
	CID     int64    // 连接ID
	UID     int64    // 用户ID
	Message *Message // 请求消息
}

// Parse 解析消息
func (r *Request) Parse(v interface{}) error {
	msg, ok := r.Message.Data.([]byte)
	if !ok {
		return copier.CopyWithOption(v, r.Message.Data, copier.Option{
			DeepCopy: true,
		})
	}

	return r.node.opts.codec.Unmarshal(msg, v)
}
