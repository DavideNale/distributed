build:
	@go build -o bin/distributed

run: build
	@./bin/distributed

test:
	@go test ./... -v
