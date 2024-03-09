.PHONY: all
all: build
FORCE: ;

.PHONY: build

build:
	go build -o bin/api-o11y-gcp cmd/api/main.go

build-linux:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags "netgo" -installsuffix netgo -o bin/api-o11y-gcp cmd/api/main.go

build-docker: 
	docker build -t api-o11y-gcp -f Dockerfile .

generate-mocks:
	@mockery --output user/mocks --dir user --all
	@mockery --output internal/telemetry/mocks --dir internal/telemetry --all

clean:
	@rm -rf user/mocks/*
	@rm -rf internal/telemetry/mocks/mocks/*

test: generate-mocks
	go test ./...

run-docker: build-docker
    docker run -d -p 8080:8080 api-o11y-gcp
