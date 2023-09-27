/**
User: cr-mao
Date: 2023/8/17 14:38
Email: crmao@qq.com
Desc: server_test.go
*/
package grpc_test

import (
	"github.com/cr-mao/loric/log"
	"github.com/cr-mao/loric/transport/grpc"
)

import (
	"testing"
)

func TestServer(t *testing.T) {
	s, err := grpc.NewServer("")
	if err != nil {
		log.Fatal(err)
	}
	log.Info(s.Addr())
	err = s.Start()
	t.Log(err)
}
