package redis_test

import (
	"context"
	"fmt"
	gate2 "github.com/cr-mao/loric/cluster/gate"
	"strconv"
	"testing"
	"time"

	"github.com/cr-mao/loric/cluster"
	"github.com/cr-mao/loric/locate/redis"
)

var locator = redis.NewLocator(
	redis.WithAddrs(
		"127.0.0.1:6379",
	),
)

func TestLocator_Set(t *testing.T) {
	for i := 1; i <= 6; i++ {
		var kind cluster.Kind

		if i%2 == 0 {
			kind = cluster.Node
		} else {
			kind = gate2.Gate
		}

		err := locator.Set(context.Background(), int64(i), kind, strconv.Itoa(i))
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestLocator_Watch(t *testing.T) {
	watcher1, err := locator.Watch(context.Background(), gate2.Gate, cluster.Node)
	if err != nil {
		t.Fatal(err)
	}

	watcher2, err := locator.Watch(context.Background(), gate2.Gate)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		for {
			events, err := watcher1.Next()
			if err != nil {
				t.Errorf("goroutine 1: %v", err)
				return
			}

			fmt.Println("goroutine 1: new event entity")

			for _, event := range events {
				t.Logf("goroutine 1: %+v", event)
			}
		}
	}()

	go func() {
		for {
			events, err := watcher2.Next()
			if err != nil {
				t.Errorf("goroutine 2: %v", err)
				return
			}

			fmt.Println("goroutine 2: new event entity")

			for _, event := range events {
				t.Logf("goroutine 2: %+v", event)
			}
		}
	}()

	time.Sleep(60 * time.Second)
}

func TestLocator_Get(t *testing.T) {
	for i := 1; i <= 6; i++ {
		insID, err := locator.Get(context.Background(), int64(i), cluster.Node)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(insID)
	}
}

func TestLocator_Rem(t *testing.T) {
	for i := 1; i <= 6; i++ {
		err := locator.Rem(context.Background(), int64(i), cluster.Node, strconv.Itoa(i))
		if err != nil {
			t.Fatal(err)
		}
	}
}
