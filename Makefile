install:
	go get

linter:
	golangci-lint run

test-coverage:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

build:
	GOOS=linux GOARCH=amd64 go build -o registry