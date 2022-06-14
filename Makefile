VERSION ?= dev

test:
	go clean -testcache && go test -race -cover ./...

install: test
	go install -ldflags "-X main.Version=${VERSION}"

build: test
	go build -ldflags "-X main.Version=${VERSION}"
