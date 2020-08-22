VSERSION ?= dev

test:
	go clean -testcache && go test -race -cover ./...

install: test
	go install -ldflags "-X main.Version=${VERSION}"

PLATFORMS := linux/amd64 darwin/amd64 linux/arm

build: test
	go build

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

release: test $(PLATFORMS)

$(PLATFORMS):
	GO111MODULE=on GOOS=$(os) GOARCH=$(arch) go build -o 'waf-tester-$(os)-$(arch)'
