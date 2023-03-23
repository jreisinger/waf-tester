test:
	go clean -testcache && go test -race -cover ./...

install: test
	go install
