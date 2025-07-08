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

format:
	golines -w .