## About

`waf-tester` runs HTTP tests against a host protected by a Web Application Firewall (WAF). The tests are defined as YAML files based on [FTW format](https://github.com/CRS-support/ftw/blob/master/docs/YAMLFormat.md) - see the `waf_tests` folder for examples. The tests are evaluated by comparing the response code (higher priority) or WAF logs against the expected values.

## Installation

Test and install (to `~/go/bin/`):

```
git clone git@github.com:jreisinger/waf-tester.git
cd waf-tester
make install
```

## Sample usage

Run some WAF tests against localhost:

```
# Run just selected tests.
waf-tester -host localhost -scheme http -tests waf_tests/generic/basic-tests.yaml

# Run all tests found in waf_tests folder. Evaluate the tests using also logs. Print overall report.
waf-tester -host localhost -scheme http -tests waf_tests/ -logs /tmp/var/log/modsec_audit.log -report
```

Consider using [waf-runner](https://github.com/jreisinger/waf-runner) to run a WAF on localhost.
