package dispatcher

import (
	"math/rand"

	"github.com/cr-mao/loric/internal/endpoint"
	"github.com/cr-mao/loric/log"
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
	//endpointWeight       map[string]int     //  id=>weight
	effectiveEndpointArr []*serviceEndpoint //  权重非0的 服务端点
}

// FindEndpoint 查询路由服务端点
func (a *abstract) FindEndpoint(insID ...string) (*endpoint.Endpoint, error) {
	if len(insID) == 0 || insID[0] == "" {
		return a.noZeroWeightRobinDispatch()
		// 暂时只有一种方式
		//switch a.dispatcher.strategy {
		//case Random:
		//	return a.randomDispatch()
		//case RoundRobin:
		//	return a.roundRobinDispatch()
		//case NotWeightZeroRoundRobin:
		//	return a.noZeroWeightRobinDispatch()
		//default:
		//	// 默认采用 轮训方式
		//	return a.roundRobinDispatch()
		//}
	}

	return a.directDispatch(insID[0])
}

// 添加服务端点
func (a *abstract) addEndpoint(insID string, ep *endpoint.Endpoint, weight int) {
	// 针对每个event 对象，route 对象， 所以这个是不会进去的
	if sep, ok := a.endpointMap[insID]; ok {
		// 不会进去的case
		sep.endpoint = ep
	} else {
		sep = &serviceEndpoint{insID: insID, endpoint: ep}
		a.endpointArr = append(a.endpointArr, sep)
		a.endpointMap[insID] = sep
		//a.endpointWeight[insID] = weight
		if weight > 0 {
			a.effectiveEndpointArr = append(a.effectiveEndpointArr, sep)
		}
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

func (a *abstract) noZeroWeightRobinDispatch() (*endpoint.Endpoint, error) {
	if len(a.endpointArr) == 0 {
		return nil, ErrNotFoundEndpoint
	}
	// 做个兜底，都是权重0的
	if len(a.effectiveEndpointArr) == 0 {
		log.Warnf("No effective endpoint")
		counter := a.counter + 1
		index := int(counter % int64(len(a.endpointArr)))
		return a.endpointArr[index].endpoint, nil
	}
	counter := a.counter + 1
	index := int(counter % int64(len(a.effectiveEndpointArr)))
	return a.effectiveEndpointArr[index].endpoint, nil
}

//nolint 随机分配
func (a *abstract) randomDispatch() (*endpoint.Endpoint, error) {
	if len(a.endpointArr) == 0 {
		return nil, ErrNotFoundEndpoint
	}
	index := rand.Int() % len(a.endpointArr)
	return a.endpointArr[index].endpoint, nil
}

//nolint 轮询分配
func (a *abstract) roundRobinDispatch() (*endpoint.Endpoint, error) {
	if len(a.endpointArr) == 0 {
		return nil, ErrNotFoundEndpoint
	}
	// 这里不提供并发保证，因为这种错误可以允许接受
	counter := a.counter + 1
	index := int(counter % int64(len(a.endpointArr)))
	return a.endpointArr[index].endpoint, nil
}
