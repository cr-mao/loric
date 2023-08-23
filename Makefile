vet:
	go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow@latest
	go vet -vettool=`which shadow` ./...
	go vet ./...

test:
	go test ./... -cover

bench:
	go test ./...  -test.bench . -test.benchmem=true

fmt:
	gofmt -w -l .

lint:
	golangci-lint cache clean
	golangci-lint run

check: fmt lint vet


install-protoc-gen-go:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28

install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.48.0





%:
	@true
