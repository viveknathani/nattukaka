build:
	go build -o ./bin/nattukaka main.go

test:
	go test -v ./...

run:
	export PORT=8080 && ./bin/nattukaka