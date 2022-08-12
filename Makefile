build:
	go build -o ./bin/nattukaka ./cmd/server
	go build -o ./bin/signup ./cmd/signup

test:
	go test -v ./...

run:
	./bin/nattukaka