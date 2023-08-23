package cluster

//  集群实例类型
const (
	Gate = "gate" // 网关服
	Node = "node" // 节点服
)

// State 集群实例状态
const (
	Work = 1 // 工作（节点正常工作，可以分配更多玩家到该节点）
	Busy = 2 // 繁忙（节点资源紧张，不建议分配更多玩家到该节点上）暂时没用 za
	Hang = 3 // 挂起（节点即将关闭，正处于资源回收中）  暂时没用
	Shut = 4 // 关闭（节点已经关闭，无法正常访问该节点）
)

// Event 事件
const (
	Connect    = 1 // 打开连接 , 暂时没用
	Reconnect  = 2 // 断线重连   node通知绑定网关,则
	Disconnect = 3 // 断开连接
)
