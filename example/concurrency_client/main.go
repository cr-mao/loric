/**
User: cr-mao
Date: 2023/8/25
Time: 15:38
Desc: 模拟app多客户端
*/
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cr-mao/loric/encoding/proto"
	"github.com/cr-mao/loric/example/internal/pb"
	"github.com/cr-mao/loric/network"
	"github.com/cr-mao/loric/network/tcp"
	"github.com/cr-mao/loric/packet"
)

func main() {
	// 并发数
	concurrency := 1000
	// 消息量
	total := 12000
	// 总共发送的消息条数
	totalSent := int64(0)
	// 总共接收的消息条数
	totalRecv := int64(0)
	// 准备消息
	wg := sync.WaitGroup{}
	client := tcp.NewClient()
	client.OnReceive(func(conn network.Conn, msg []byte) {
		atomic.AddInt64(&totalRecv, 1)
		message, uErr := packet.Unpack(msg)
		if uErr != nil {
			fmt.Println(uErr)
		}
		if message.Seq != 0 {
			fmt.Println("seq error")
		}
		if message.Route != 0 {
			fmt.Println("Route error")
		}
		wg.Done()
	})

	wg.Add(total)

	chMsg := make(chan struct{}, total)

	// 准备连接
	conns := make([]network.Conn, concurrency)
	for i := 0; i < concurrency; i++ {
		conn, dErr := client.Dial()
		if dErr != nil {
			fmt.Println("connect failed", i, dErr)
			i--
			continue
		}

		conns[i] = conn
		time.Sleep(time.Millisecond * 2)
	}

	// 发送消息
	for _, conn := range conns {
		go func(conn network.Conn) {
			defer conn.Close(true)

			for {
				select {
				case _, ok := <-chMsg:
					if !ok {
						return
					}
					msg1, _ := proto.Marshal(&pb.LoginReq{
						Token: "cr-mao",
					})
					msg, _ := packet.Pack(&packet.Message{
						Seq:    0,
						Route:  int32(pb.Route_Login),
						Buffer: msg1,
					})
					if err := conn.Push(msg); err != nil {
						fmt.Println(err)
						return
					}
					atomic.AddInt64(&totalSent, 1)
				}
			}
		}(conn)

	}

	startTime := time.Now().UnixNano()

	for i := 0; i < total; i++ {
		chMsg <- struct{}{}
	}

	wg.Wait()
	close(chMsg)

	totalTime := float64(time.Now().UnixNano()-startTime) / float64(time.Second)

	/*
		server               : tcp
		concurrency          : 1000
		latency              : 66.533924s
		sent     requests    : 1000000
		received requests    : 1000000
		throughput  (TPS)    : 15029
	*/

	fmt.Printf("server               : %s\n", "tcp")
	fmt.Printf("concurrency          : %d\n", concurrency)
	fmt.Printf("latency              : %fs\n", totalTime)
	fmt.Printf("sent     requests    : %d\n", totalSent)
	fmt.Printf("received requests    : %d\n", totalRecv)
	fmt.Printf("throughput  (TPS)    : %d\n", int64(float64(totalRecv)/totalTime))

}
