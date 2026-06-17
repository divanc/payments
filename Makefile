.PHONY: run test build tidy

run:
	go run .

test:
	go test ./...

build:
	go build -o payments .

tidy:
	go mod tidy
