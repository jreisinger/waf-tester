[![Build Status](https://travis-ci.org/jreisinger/waf-tester.svg?branch=master)](https://travis-ci.org/jreisinger/waf-tester)

## About

waf-tester runs tests against a URL protected by a Web Application Firewall (WAF). The tests are HTTP requests defined in YAML format based on [FTW format](https://github.com/CRS-support/ftw/blob/master/docs/YAMLFormat.md). Use '-template' to see how they look like.

The tests are evaluated by comparing the HTTP response status or WAF logs against the expected values defined in tests. If both 'status' and 'log_contains' are defined in a test only status is evaluated. If '-logs' is not used tests containing only 'log_contains' are skipped.

## Installation

`make install`

## Sample usage

Run some WAF tests against localhost:

```
# Generate tests and run them against localhost.
waf-tester -template > tests.yaml
waf-tester -verbose

# Run tests from waf_tests folder and evaluate also logs (NOTE: -logs is kind of experimental).
waf-tester -tests waf_tests/ -logs /tmp/var/log/modsec_audit.log
```

Consider using [waf-runner](https://github.com/jreisinger/waf-runner) to run a WAF on localhost.

## Development

```
vim main.go
make build      # see dist folder
make install    # install into ~/go/bin
```

Check test coverage:

```
go test -coverprofile cover.out ./...
go tool cover -html=cover.out
```

If the code is slow [profile](https://blog.golang.org/pprof) it.
