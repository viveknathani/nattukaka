build:
	go build -o ./bin/nattukaka ./cmd/cli
	go build -o ./bin/nattukaka-server ./cmd/server
	
test:
	go test -v ./...

run:
	./bin/nattukaka-server