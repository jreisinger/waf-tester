*NOTE*: This is still experimental.

waf-tester runs HTTP tests against hosts protected by a WAF. The tests are defined as YAML files based on [FTW format](https://github.com/CRS-support/ftw/blob/master/docs/YAMLFormat.md).

Install and run:

```
go build
./waf-tester -h
./waf-tester -o tests/ok-tests-nginx-modsecurity-2020-01-22.txt -s
```
