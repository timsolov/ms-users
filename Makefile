BUF_VERSION:=1.0.0-rc9

.PHONY: build
build:
	buf generate
	go build -o service ./cmd/

generate:
	docker run -v $$(pwd):/src -w /src --rm bufbuild/buf:$(BUF_VERSION) generate

lint:
	docker run -v $$(pwd):/src -w /src --rm bufbuild/buf:$(BUF_VERSION) lint
	docker run -v $$(pwd):/src -w /src --rm bufbuild/buf:$(BUF_VERSION) breaking --against 'https://github.com/timsolov/ms-users.git#branch=master'
