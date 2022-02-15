BUF_VERSION:=1.1.0

.DEFAULT_GOAL := default

.PHONY: default
default: generate build

.PHONY: build
build:
	go build -o service ./cmd/

.PHONY: generate
generate:
	docker run -v $$(pwd):/src -w /src --rm bufbuild/buf:$(BUF_VERSION) generate

.PHONY: lint
lint:
	docker run -v $$(pwd):/src -w /src --rm bufbuild/buf:$(BUF_VERSION) lint
	docker run -v $$(pwd):/src -w /src --rm bufbuild/buf:$(BUF_VERSION) breaking --against 'https://github.com/timsolov/ms-users.git#branch=main'

.PHONY: tools
tools:
	go install github.com/golang/mock/mockgen
	go install github.com/golangci/golangci-lint/cmd/golangci-lint
	go install github.com/bufbuild/buf/cmd/buf
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	go install github.com/fzipp/gocyclo/cmd/gocyclo
	go install github.com/go-critic/go-critic/cmd/gocritic