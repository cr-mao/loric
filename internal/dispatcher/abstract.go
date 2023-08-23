package dispatcher

import (
	"math/rand"
	"sync/atomic"

	"github.com/cr-mao/loric/internal/endpoint"
)

type serviceEndpoint struct {
	insID    string
	endpoint *endpoint.Endpoint
}

type abstract struct {
	counter     int64
	dispatcher  *Dispatcher
	endpointMap map[string]*serviceEndpoint
	endpointArr []*serviceEndpoint
}

// FindEndpoint 查询路由服务端点
func (a *abstract) FindEndpoint(insID ...string) (*endpoint.Endpoint, error) {
	if len(insID) == 0 || insID[0] == "" {
		switch a.dispatcher.strategy {
		case Random:
			return a.randomDispatch()
		case RoundRobin:
			return a.roundRobinDispatch()
		default:
			// 默认采用 轮训方式
			return a.roundRobinDispatch()
		}
	}

	return a.directDispatch(insID[0])
}

// 添加服务端点
func (a *abstract) addEndpoint(insID string, ep *endpoint.Endpoint) {
	if sep, ok := a.endpointMap[insID]; ok {
		sep.endpoint = ep
	} else {
		sep = &serviceEndpoint{insID: insID, endpoint: ep}
		a.endpointArr = append(a.endpointArr, sep)
		a.endpointMap[insID] = sep
	}
}

// 直接分配
func (a *abstract) directDispatch(insID string) (*endpoint.Endpoint, error) {
	sep, ok := a.endpointMap[insID]
	if !ok {
		return nil, ErrNotFoundEndpoint
	}

	return sep.endpoint, nil
}

// 随机分配
func (a *abstract) randomDispatch() (*endpoint.Endpoint, error) {
	if len(a.endpointArr) == 0 {
		return nil, ErrNotFoundEndpoint
	}
	index := rand.Int() % len(a.endpointArr)
	return a.endpointArr[index].endpoint, nil
}

// 轮询分配
func (a *abstract) roundRobinDispatch() (*endpoint.Endpoint, error) {
	if len(a.endpointArr) == 0 {
		return nil, ErrNotFoundEndpoint
	}
	counter := atomic.AddInt64(&a.counter, 1)
	index := int(counter % int64(len(a.endpointArr)))
	return a.endpointArr[index].endpoint, nil
}
