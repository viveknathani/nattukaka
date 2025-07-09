PROTO_SRC=proto/player.proto
PROTO_OUT=gen

ifneq (,$(wildcard .env))
	include .env
endif

build:
	go build -o ./bin/main main.go

test:
	go test -v ./...

run-dev:
	export ENV=development && make build && ./bin/main

run-production:
	export ENV=production && ./bin/main

proto:
	protoc --go_out=$(PROTO_OUT) --go-grpc_out=$(PROTO_OUT) $(PROTO_SRC)

format:
	golines -w .
