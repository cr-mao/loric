package dispatcher_test

import (
	"github.com/cr-mao/loric/cluster"
	"github.com/cr-mao/loric/internal/dispatcher"
	"github.com/cr-mao/loric/internal/endpoint"
	"github.com/cr-mao/loric/registry"
	"math/rand"
	"testing"
	"time"
)

func TestDispatcher_ReplaceServices(t *testing.T) {
	var (
		instance1 = &registry.ServiceInstance{
			ID:       "xc",
			Name:     "gate-3",
			Kind:     cluster.Node,
			Alias:    "gate-3",
			State:    cluster.Work,
			Endpoint: endpoint.NewEndpoint("grpc", "127.0.0.1:8003", false).String(),
			Routes: []registry.Route{{
				ID:       2,
				Stateful: false,
			}, {
				ID:       3,
				Stateful: false,
			}, {
				ID:       4,
				Stateful: true,
			}},
		}
		instance2 = &registry.ServiceInstance{
			ID:       "xa",
			Name:     "gate-1",
			Kind:     cluster.Node,
			Alias:    "gate-1",
			State:    cluster.Work,
			Endpoint: endpoint.NewEndpoint("grpc", "127.0.0.1:8001", false).String(),
			Routes: []registry.Route{{
				ID:       1,
				Stateful: false,
			}, {
				ID:       2,
				Stateful: false,
			}, {
				ID:       3,
				Stateful: false,
			}, {
				ID:       4,
				Stateful: true,
			}},
		}
		instance3 = &registry.ServiceInstance{
			ID:       "xb",
			Name:     "gate-2",
			Kind:     cluster.Node,
			Alias:    "gate-2",
			State:    cluster.Work,
			Endpoint: endpoint.NewEndpoint("grpc", "127.0.0.1:8002", false).String(),
			Events:   []int32{cluster.Disconnect},
			Routes: []registry.Route{{
				ID:       1,
				Stateful: false,
			}, {
				ID:       2,
				Stateful: false,
			}},
		}
	)

	d := dispatcher.NewDispatcher(dispatcher.RoundRobin)

	d.ReplaceServices(instance1, instance2, instance3)

	event, err := d.FindEvent(cluster.Disconnect)
	if err != nil {
		t.Errorf("find event failed: %v", err)
	} else {
		t.Log(event.FindEndpoint())
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

//Benchmark_StrangerRoundRobin-8   	  374792	      3064 ns/op
func Benchmark_StrangerRoundRobin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var (
			instance1 = &registry.ServiceInstance{
				ID:       "xc",
				Name:     "gate-3",
				Kind:     cluster.Node,
				Alias:    "gate-3",
				State:    cluster.Work,
				Endpoint: endpoint.NewEndpoint("grpc", "127.0.0.1:8003", false).String(),
				Routes: []registry.Route{{
					ID:       2,
					Stateful: false,
				}, {
					ID:       3,
					Stateful: false,
				}, {
					ID:       4,
					Stateful: true,
				}},
			}
		)

		d := dispatcher.NewDispatcher(dispatcher.RoundRobin)

		d.ReplaceServices(instance1)

		router, _ := d.FindRoute(2)
		router.FindEndpoint()

	}

}

//Benchmark_StrangerRandom-8   	  388052	      3087 ns/op
func Benchmark_StrangerRandom(b *testing.B) {

	for i := 0; i < b.N; i++ {

		var (
			instance1 = &registry.ServiceInstance{
				ID:       "xc",
				Name:     "gate-3",
				Kind:     cluster.Node,
				Alias:    "gate-3",
				State:    cluster.Work,
				Endpoint: endpoint.NewEndpoint("grpc", "127.0.0.1:8003", false).String(),
				Routes: []registry.Route{{
					ID:       2,
					Stateful: false,
				}, {
					ID:       3,
					Stateful: false,
				}, {
					ID:       4,
					Stateful: true,
				}},
			}
		)
		d := dispatcher.NewDispatcher(dispatcher.Random)
		d.ReplaceServices(instance1)
		router, _ := d.FindRoute(2)
		router.FindEndpoint()
	}

}
