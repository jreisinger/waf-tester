waf-tester runs HTTP tests against hosts protected by a WAF. The tests are defined as YAML files based on [FTW format](https://github.com/CRS-support/ftw/blob/master/docs/YAMLFormat.md).

Install and run:

```
go build
./waf-tester -h
./waf-tester localhost
```

To do:

* [x] parse [FTW YAML](https://github.com/CRS-support/ftw/blob/master/docs/YAMLFormat.md) so I can use https://github.com/SpiderLabs/owasp-modsecurity-crs/tree/v3.3/dev/tests (https://medium.com/@kenanbek/golang-how-to-parse-yaml-file-31b78141bda7)
