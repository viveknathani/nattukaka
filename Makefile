build:
	go build -o ./bin/nattukaka ./cmd/server
	go build -o ./bin/signup ./cmd/signup
	go build -o ./bin/reporter ./cmd/reporter
test:
	go test -v ./...

run:
	./bin/nattukaka