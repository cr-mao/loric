package tcp_test

import (
	"net/http"
	_ "net/http/pprof"
	"syscall"
	"testing"

	"github.com/cr-mao/loric/log"
	"github.com/cr-mao/loric/network"
	"github.com/cr-mao/loric/network/tcp"
	"github.com/cr-mao/loric/packet"
)

func setLimit() {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
}

func TestServer(t *testing.T) {
	t.SkipNow()
	setLimit()
	server := tcp.NewServer()
	server.OnStart(func() {
		t.Logf("server", server.Addr(), " is started")
	})
	server.OnConnect(func(conn network.Conn) {
		t.Logf("connection is opened, connection id: %d", conn.ID())
	})
	server.OnDisconnect(func(conn network.Conn) {
		t.Logf("connection is closed, connection id: %d", conn.ID())
	})
	server.OnReceive(func(conn network.Conn, msg []byte) {
		message, err := packet.Unpack(msg)
		if err != nil {
			t.Error(err)
			return
		}
		t.Logf("receive msg from client, connection id: %d, seq: %d, route: %d, msg: %s", conn.ID(), message.Seq, message.Route, string(message.Buffer))
		msg, err = packet.Pack(&packet.Message{
			Seq:    1,
			Route:  1,
			Buffer: []byte("I'm fine~~"),
		})
		if err != nil {
			t.Error(err)
			return
		}

		if err = conn.Push(msg); err != nil {
			t.Error(err)
		}
	})

	if err := server.Start(); err != nil {
		t.Fatal(err)
	}

	go func() {
		err := http.ListenAndServe(":8089", nil)
		if err != nil {
			log.Errorf("pprof server start failed: %v", err)
		}
	}()

	select {}
}

func TestServer_Benchmark(t *testing.T) {
	t.SkipNow()
	setLimit()

	server := tcp.NewServer()
	server.OnStart(func() {
		t.Logf("server is started")
	})
	server.OnReceive(func(conn network.Conn, msg []byte) {
		_, err := packet.Unpack(msg)
		if err != nil {
			t.Error(err)
			return
		}

		msg, err = packet.Pack(&packet.Message{
			Seq:    1,
			Route:  1,
			Buffer: []byte("I'm fine~~"),
		})
		if err != nil {
			t.Error(err)
			return
		}

		if err = conn.Push(msg); err != nil {
			t.Error(err)
		}
	})

	if err := server.Start(); err != nil {
		t.Fatal(err)
	}

	go func() {
		err := http.ListenAndServe(":8089", nil)
		if err != nil {
			log.Errorf("pprof server start failed: %v", err)
		}
	}()

	select {}
}
