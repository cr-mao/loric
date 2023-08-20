package node

import (
	"sync"

	"github.com/cr-mao/loric/cluster"
	"github.com/cr-mao/loric/log"
)

type EventHandler func(event *Event)

type Event struct {
	Proxy *Proxy
	Event int8
	GID   string
	CID   int64
	UID   int64
}

type Events struct {
	node    *Node
	events  map[int8]EventHandler
	evtPool sync.Pool
	evtChan chan *Event
}

func newEvents(node *Node) *Events {
	return &Events{
		node:    node,
		events:  make(map[int8]EventHandler, 3),
		evtPool: sync.Pool{New: func() interface{} { return &Event{Proxy: node.proxy} }},
		evtChan: make(chan *Event, 4096),
	}
}

// 触发事件
func (e *Events) trigger(event int8, gid string, cid, uid int64) {
	evt := e.evtPool.Get().(*Event)
	evt.Event = event
	evt.GID = gid
	evt.CID = cid
	evt.UID = uid
	e.evtChan <- evt
}

func (e *Events) receive() <-chan *Event {
	return e.evtChan
}

func (e *Events) close() {
	close(e.evtChan)
}

func (e *Events) handle(evt *Event) {
	defer e.evtPool.Put(evt)

	handler, ok := e.events[evt.Event]
	if !ok {
		return
	}

	handler(evt)
}

// AddEventHandler 添加事件处理器
func (e *Events) AddEventHandler(event int8, handler EventHandler) {
	if e.node.getState() != cluster.Shut {
		log.Warnf("the node server is working, can't add Event handler")
		return
	}

	e.events[event] = handler
}
