`waf-tester` runs HTTP tests against a host protected by a Web Application Firewall (WAF). The tests are defined as YAML files based on [FTW format](https://github.com/CRS-support/ftw/blob/master/docs/YAMLFormat.md) (see the `tests*` folders). The tests are evaluated by checking the response code (higher priority) or WAF logs.

Test and install (to `~/go/bin/`):

```
make install
```

Sample usage:

```
# Start a testing WAF and a testing web server on localhost.
./run-waf -s nginx/modsecurity

# In a different terminal window run some WAF tests against localhost.
waf-tester -host localhost -scheme http -tests tests/
```

Test and build (cross-compile for multiple platforms):

```
make release
```
