/**
User: cr-mao
Date: 2023/8/17 14:38
Email: crmao@qq.com
Desc: server_test.go
*/
package grpc_test

import (
	"github.com/cr-mao/loric/transport/grpc"
)

import (
	"context"
	"testing"
)

func TestServer(t *testing.T) {
	t.SkipNow()
	s := grpc.NewServer(grpc.WithAddress("0.0.0.0:9000"))
	ctx := context.Background()
	err := s.Start(ctx)
	t.Log(err)
}
