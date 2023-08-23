package node

import (
	"github.com/cr-mao/loric/sugar"
	"sync"

	"github.com/cr-mao/loric/cluster"
	"github.com/cr-mao/loric/log"
)

type EventHandler func(event *Event)

type Event struct {
	Proxy *Proxy
	Event int32
	GID   string
	CID   int64
	UID   int64
}

type Events struct {
	node    *Node
	events  map[int32]EventHandler
	evtPool sync.Pool
}

func newEvents(node *Node) *Events {
	return &Events{
		node:    node,
		events:  make(map[int32]EventHandler, 3),
		evtPool: sync.Pool{New: func() interface{} { return &Event{Proxy: node.proxy} }},
	}
}

// AddEventHandler 添加事件处理器
func (e *Events) AddEventHandler(event int32, handler EventHandler) {
	if e.node.getState() != cluster.Shut {
		log.Warnf("the node server is working, can't add Event handler")
		return
	}
	e.events[event] = handler
	log.Debugf("add event %s --> handler: %s", cluster.EventNames[event], sugar.NameOfFunction(handler))
}
