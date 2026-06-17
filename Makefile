.PHONY: run test e2e build tidy

run:
	go run .

test:
	go test ./...

e2e:
	go test -tags e2e ./e2e/

build:
	go build -o payments .

tidy:
	go mod tidy
