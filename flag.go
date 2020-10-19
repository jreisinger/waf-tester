package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

// Flags are all the available CLI flags (options).
type Flags struct {
	Verbose     bool
	TestsPath   string
	LogsPath    string
	URL         string
	Title       string
	Template    bool
	Version     bool
	Rate        int
	Print       string
	Concurrency int
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

// ParseFlags validates the flags and parses them into Flags.
func ParseFlags() (Flags, error) {
	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	URL := f.String("url", "http://localhost", "`URL` to test")
	TestsPath := f.String("tests", "", "`DIR|FILE` containing tests")
	Verbose := f.Bool("verbose", false, "print more info about tests")
	LogsPath := f.String("logs", "", "evaluate logs from `FILE|API` (e.g. modsec_audit.log or https://loki.example.com)")
	Title := f.String("title", "", "execute only test with `TITLE`")
	Template := f.Bool("template", false, "print tests template and exit")
	Version := f.Bool("version", false, "version")
	Rate := f.Int("rate", 100, "maximum number of HTTP requests per second")
	Print := f.String("print", "", "print info about tests with status `FAIL|OK|ERR`")
	Concurrency := f.Int("conc", 10, "maximum number of concurrent HTTP requests")

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
		Title:       stringValue(Title),
		Template:    boolValue(Template),
		Version:     boolValue(Version),
		Rate:        intValue(Rate),
		Print:       stringValue(Print),
		Concurrency: intValue(Concurrency),
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
