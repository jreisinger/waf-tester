## About

`waf-tester` runs HTTP tests against a host protected by a Web Application Firewall (WAF). The tests are defined as YAML files based on [FTW format](https://github.com/CRS-support/ftw/blob/master/docs/YAMLFormat.md) (see the `tests*` folders). The tests are evaluated by checking the response code (higher priority) or WAF logs.

`run-waf` is a helper tool for running WAFs locally as Docker containers for testing purposes.

## Installation

Test and install `waf-tester` (to `~/go/bin/`):

```
make install
```

## Sample usage

Start a testing WAF (and a testing web server) on localhost:

```
./run-waf -s nginx/modsecurity
```

In a different terminal window run some WAF tests against localhost:

```
waf-tester -host localhost -scheme http -tests tests/
```

Obviously, you can run `waf-tester` against any host not just against our testing WAF container on localhost. See this [blog post](https://jreisinger.github.io/blog2/posts/waf-tester/) for more.