`waf-tester` runs HTTP tests against a host protected by a Web Application Firewall (WAF). The tests are defined as YAML files based on [FTW format](https://github.com/CRS-support/ftw/blob/master/docs/YAMLFormat.md). You can write you own tests (see the `tests` folder) or you can re-use existing tests (see the `tests-crs` folder). The tests are evaluated by checking the HTTP response status code (higher priority) or WAF logs.

Test and install (to `~/go/bin/`):

```
make install
```

Sample usage:

```
# Start a testing WAF and a testing web server on localhost.
./run-waf -s modsecurity

# In a different terminal window run some WAF tests against localhost.
waf-tester -host localhost -scheme http -tests tests
```

Test and build (cross-compile for multiple platforms):

```
make release
```
