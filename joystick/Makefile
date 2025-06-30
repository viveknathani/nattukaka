ifneq (,$(wildcard .env))
	include .env
endif

build:
	go build -o ./bin/main main.go

test:
	go test -v ./...


migrations-create:
	migrate create -ext sql -dir ./database/migrations $(name)

migrations-up:
	migrate -verbose -path ./database/migrations -database $(DATABASE_URL) up

migrations-down:
	migrate -verbose -path ./database/migrations -database $(DATABASE_URL) down 1

run-dev:
	export ENV=development && make build && ./bin/main

run-production:
	export ENV=production && ./bin/main

format:
	golines -w .