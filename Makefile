build:
	go build -o ./bin/nattukaka main.go

test:
	go test -v ./...

run:
	./bin/nattukaka