FROM golang:1.15

WORKDIR /src
COPY . .
RUN go test -race -cover ./...

RUN mkdir /releases
ARG VERSION=dev

RUN GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$VERSION" -o /releases/waf-tester
RUN tar -czvf /releases/waf-tester_linux_amd64.tar.gz -C /releases/ waf-tester
RUN rm /releases/waf-tester

RUN GOOS=linux GOARCH=arm go build -ldflags "-X main.Version=$VERSION" -o /releases/waf-tester
RUN tar -czvf /releases/waf-tester_linux_arm.tar.gz -C /releases/ waf-tester
RUN rm /releases/waf-tester

RUN GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.Version=$VERSION" -o /releases/waf-tester
RUN tar -czvf /releases/waf-tester_darwin_amd64.tar.gz -C /releases/ waf-tester
RUN rm /releases/waf-tester