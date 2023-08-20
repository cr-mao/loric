package dispatcher

type Event struct {
	abstract
	event int32 // 路由ID
}

func newEvent(dispatcher *Dispatcher, event int32) *Event {
	return &Event{
		abstract: abstract{
			dispatcher:  dispatcher,
			endpointMap: make(map[string]*serviceEndpoint),
			endpointArr: make([]*serviceEndpoint, 0),
		},
		event: event,
	}
}

// Event 获取事件
func (e *Event) Event() int32 {
	return e.event
}
