VSERSION ?= dev

test:
	go clean -testcache && go test -race -cover ./...

install: test
	go install -ldflags "-X main.Version=${VERSION}"

build: test
	go build

release:
	docker build --build-arg version=${VERSION} -t waf-tester-releases -f Releases.Dockerfile .
	docker create -ti --name waf-tester-releases waf-tester-releases sh
	test -d releases || mkdir releases
	docker cp waf-tester-releases:/releases/waf-tester_linux_amd64.tar.gz releases/
	docker cp waf-tester-releases:/releases/waf-tester_linux_arm.tar.gz releases/
	docker cp waf-tester-releases:/releases/waf-tester_darwin_amd64.tar.gz releases/
	docker rm waf-tester-releases
	docker rmi waf-tester-releases
