*NOTE*: This is still experimental.

waf-tester runs HTTP tests against hosts protected by a WAF. The tests are defined as YAML files based on [FTW format](https://github.com/CRS-support/ftw/blob/master/docs/YAMLFormat.md). This project was inspired by [FTW](https://github.com/CRS-support/ftw).

Test and install (to `~/go/bin/`):

```
make install
```

Test and build (cross-compile for multiple platforms):

```
make release
```