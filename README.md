## About

`waf-tester` runs HTTP tests against a host protected by a Web Application Firewall (WAF). The tests are defined as YAML files based on [FTW format](https://github.com/CRS-support/ftw/blob/master/docs/YAMLFormat.md) - see the `tests` folder for examples. The tests are evaluated by checking the response code (higher priority) or WAF logs against the exptected values. See [this blog post](https://jreisinger.github.io/blog2/posts/working-with-waf-containers/) for more.

## Installation

Test and install (to `~/go/bin/`):

```
make install
```

## Sample usage

Run some WAF tests against localhost:

```
waf-tester -host localhost -scheme http -tests tests/basic-tests.yaml
waf-tester -host localhost -scheme http -tests tests/ -logs /tmp/var/log/modsec_audit.log -report
```

Consider using [waf-runner](https://github.com/jreisinger/waf-runner) to run a WAF on localhost.
