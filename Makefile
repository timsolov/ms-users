BUF_VERSION:=1.1.0

.DEFAULT_GOAL := default

.PHONY: default
default: generate build

.PHONY: build
build:
	go build -o service ./cmd/

.PHONY: gen
gen:
	buf generate --template api/proto/buf.gen.yaml api/proto

.PHONY: lint
lint:
	buf lint api/proto

.PHONY: tools
tools:
	go install github.com/golang/mock/mockgen@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	go install github.com/go-critic/go-critic/cmd/gocritic@latest