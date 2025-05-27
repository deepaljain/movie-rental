build:
	go build ./cmd

run:
	go run ./cmd/main.go

test:
	go test -v ./...

migrateup:
	migrate -path ./migrations -database "postgres://deepaljain:postgres@localhost:5432/movierental?sslmode=disable" up

migratedown:
	migrate -path ./migrations -database "postgres://deepaljain:postgres@localhost:5432/movierental?sslmode=disable" down