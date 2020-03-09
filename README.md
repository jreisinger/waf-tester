*NOTE*: This is still experimental.

`waf-tester` runs HTTP tests against a host protected by a Web Application Firewall (WAF). The tests are defined as YAML files based on [FTW format](https://github.com/CRS-support/ftw/blob/master/docs/YAMLFormat.md).

Test and install (to `~/go/bin/`):

```
make install
```

Test and build (cross-compile for multiple platforms):

```
make release
```
