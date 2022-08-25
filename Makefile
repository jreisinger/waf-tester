test:
	go clean -testcache && go test -race -cover ./...

build: test
	goreleaser build --single-target --rm-dist

install: build
	cp dist/waf-tester_darwin_amd64_v1/waf-tester ~/go/bin/
