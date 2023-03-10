build:
	@go build -o bin/netassertv2-l4-client main.go

build-race:
	@go build -race -o bin/netassertv2-l4-client main.go

lint:
	golangci-lint run ./...

clean:
	@rm -rf bin/netassertv2-l4-client

docker-build:
	docker build -f Dockerfile \
	--no-cache \
    --tag local/netassertv2-l4-client:dev .
