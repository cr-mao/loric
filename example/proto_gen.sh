#!/bin/sh

find internal/pb -type f -name "*.pb.go" -delete

absoulte_path=`pwd`

# source_relative  模式， 和go_out， go-grpc_out 路径一致
protoc --proto_path=internal/pb --go_out=$absoulte_path/internal/pb --go_opt=paths=source_relative  game.proto










