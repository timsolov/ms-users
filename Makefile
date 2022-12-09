ifeq ($(shell git tag --contains HEAD),)
  VERSION := $(shell git rev-parse --short HEAD)
else
  VERSION := $(shell git tag --contains HEAD)
endif

BUILDTIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GOLDFLAGS += -X ms-users/app/conf.Version=$(VERSION)
GOLDFLAGS += -X ms-users/app/conf.Buildtime=$(BUILDTIME)
GOALFLAGS = -ldflags "$(GOLDFLAGS)"

.DEFAULT_GOAL := default

.PHONY: default
default: gen lint build

.PHONY: build
build:
	go build -o service $(GOALFLAGS) ./cmd/

.PHONY: gen
gen:
	clang-format -i api/proto/users/v1/users.proto
	buf generate --template api/proto/buf.gen.yaml api/proto
	buf generate --template api/proto/buf.gen.tagger.yaml api/proto
	go generate ./...

.PHONY: lint
lint:
	buf lint api/proto
	golangci-lint -c .golangci.yml run ./...

.PHONY: sql
sql: # make sql q="select * from table"
	docker exec -it ${DB_CONTAINER} psql --host=${DB_HOST} --dbname=${DB_NAME} --username=${DB_USER} --command="${q}"

.PHONY: tools
tools:
	go install github.com/golang/mock/mockgen@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	go install github.com/go-critic/go-critic/cmd/gocritic@v0.6.3
	go install github.com/srikrsna/protoc-gen-gotag@v0.6.2