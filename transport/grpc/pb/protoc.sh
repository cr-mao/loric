#!/bin/bash

# 在当前目录执行后， 把他移动到外面即可
protoc --go_out=. --go-grpc_out=. *.proto

