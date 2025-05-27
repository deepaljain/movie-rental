build:
	go build ./cmd

run:
	go run ./cmd/main.go

test:
	go test -v ./...