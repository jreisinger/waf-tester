package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

// Flags are all the available CLI flags (options).
type Flags struct {
	Verbose     bool
	TestsPath   string
	LogsPath    string
	URL         string
	Execute     arrayFlags
	NoExecute   arrayFlags
	Header      arrayFlags
	Template    bool
	Version     bool
	Print       string
	Rate        int
	Concurrency int
	Timeout     int
}

const usage = `ABOUT

waf-tester runs tests against a URL protected by a Web Application Firewall
(WAF). The tests are HTTP requests defined in YAML format. Use '-template' to
see how they look like.

EXAMPLE

# Generate and run tests.
waf-tester -template > tests.yaml
waf-tester -tests tests.yaml

OPTIONS

`

type arrayFlags []string

func (a *arrayFlags) String() string {
	return fmt.Sprintf("%s", *a)
}
func (a *arrayFlags) Set(value string) error {
	values := strings.Split(value, ",")
	for i := range values {
		values[i] = strings.TrimSpace(values[i])
	}
	*a = append(*a, values...)
	return nil
}

// ParseFlags validates the flags and parses them into Flags.
func ParseFlags() (Flags, error) {
	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	var Execute arrayFlags
	var NoExecute arrayFlags
	var Header arrayFlags

	URL := f.String("url", "http://localhost", "`URL` to test")
	TestsPath := f.String("tests", "waf_tests", "`DIR|FILE` containing tests")
	Verbose := f.Bool("verbose", false, "print more info about tests")
	LogsPath := f.String("logs", "", "evaluate logs from `FILE|API` (e.g. modsec_audit.log or https://loki.example.com)")
	f.Var(&Execute, "exec", "execute only tests with `TITLE|TAG[,...]`")
	f.Var(&NoExecute, "noexec", "don't execute tests with `TITLE|TAG[,...]`")
	f.Var(&Header, "header", "add `KEY:VALUE[,...]` header to all requests (waf-tester-id is always added)")
	Template := f.Bool("template", false, "print tests template and exit")
	Version := f.Bool("version", false, "version")
	Print := f.String("print", "", "print info about tests with status `FAIL|OK|ERR`")
	Rate := f.Int("rate", 200, "maximum number of HTTP requests per second")
	Concurrency := f.Int("conc", 20, "maximum number of concurrent HTTP requests")
	Timeout := f.Int("timeout", 10, "timeout in seconds for HTTP requests, 0 means no timeout")

	f.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), usage)
		f.PrintDefaults()
	}

	err := f.Parse(os.Args[1:])
	if err != nil {
		return Flags{}, err
	}

	flags := Flags{
		URL:         stringValue(URL),
		TestsPath:   stringValue(TestsPath),
		Verbose:     boolValue(Verbose),
		LogsPath:    stringValue(LogsPath),
		Execute:     Execute,
		NoExecute:   NoExecute,
		Header:      Header,
		Template:    boolValue(Template),
		Version:     boolValue(Version),
		Print:       stringValue(Print),
		Rate:        intValue(Rate),
		Concurrency: intValue(Concurrency),
		Timeout:     intValue(Timeout),
	}

	err = flags.validate()
	return flags, err
}

func (f Flags) validate() error {
	if !f.Template && !f.Version {
		if f.TestsPath == "" {
			return errors.New("tests cannot be empty")
		}
	}
	if f.Print != "" &&
		!(f.Print == "FAIL" || f.Print == "OK" || f.Print == "ERR") {
		return errors.New("status must be FAIL, OK or ERR")
	}
	if f.Rate < 1 {
		return errors.New("rate must be > 1")
	}
	if f.Concurrency < 1 {
		return errors.New("conc must be > 1")
	}
	if f.Timeout < 0 {
		return errors.New("timeout must be > 0")
	}
	return nil
}

func stringValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func boolValue(v *bool) bool {
	if !*v {
		return false
	}
	return *v
}

func intValue(v *int) int {
	if *v == 0 {
		return 0
	}
	return *v
}
