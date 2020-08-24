[![Build Status](https://travis-ci.org/jreisinger/waf-tester.svg?branch=master)](https://travis-ci.org/jreisinger/waf-tester)

## About

`waf-tester` runs HTTP tests against a host protected by a Web Application Firewall (WAF). The tests are defined as YAML files based on [FTW format](https://github.com/CRS-support/ftw/blob/master/docs/YAMLFormat.md). Run `waf-tester -template` or see the `waf_tests` folder for examples. The tests are evaluated by comparing the response code (higher priority) or WAF logs against the expected values.

## Installation

Download the latest [release](https://github.com/jreisinger/waf-tester/releases) for your operating system and architecture.

## Sample usage

Run some WAF tests against localhost:

```
# Generate and run tests.
waf-tester -template > tests.yaml
waf-tester -host localhost -scheme http -tests tests.yaml 

# Run all tests found in waf_tests folder. Print overall report.
waf-tester -host localhost -scheme http -tests waf_tests/ -report

# Evaluate also tests using only logs. NOTE: -logs is experimental and somehow slow.
waf-tester -host localhost -scheme http -tests waf_tests/ -logs /tmp/var/log/modsec_audit.log
```

Consider using [waf-runner](https://github.com/jreisinger/waf-runner) to run a WAF on localhost.

## Development

```
vim main.go
make install # version defaults to "dev" if VERSION envvar is not set

make release # you'll find releases in releases/ directory
```

Builds are done inside Docker container. Once you push to GitHub Travis will
try and build a release for you and publish it on GitHub.

Check test coverage:

```
go test -coverprofile cover.out ./...
go tool cover -html=cover.out
```
