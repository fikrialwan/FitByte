# application commands
build:
	go build -o bin/main cmd/app/main.go

run:
	go run cmd/app/main.go

build-run: build
	./bin/main

test:
	go test -v ./tests

# database commands
migrate:
	go run cmd/migrate/main.go --migrate
