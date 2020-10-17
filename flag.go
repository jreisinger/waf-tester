package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

// Flags are all the available CLI flags (options).
type Flags struct {
	Verbose   bool
	TestsPath string
	LogsPath  string
	URL       string
	Title     string
	Template  bool
	Version   bool
	RPS       int
	Print     string
}

const usage = `ABOUT

waf-tester runs tests against a URL protected by a Web Application Firewall
(WAF). The tests are HTTP requests defined in YAML format. Use '-template' to
see how they look like.

EXAMPLE

# Generate and run tests.
waf-tester -template > tests.yaml
waf-tester -url http://localhost -tests tests.yaml

OPTIONS

`

// ParseFlags validates the flags and parses them into Flags.
func ParseFlags() (Flags, error) {
	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	URL := f.String("url", "", "`URL` to test (e.g. https://example.com)")
	TestsPath := f.String("tests", "", "`DIR|FILE` containing tests")
	Verbose := f.Bool("verbose", false, "print more info about tests")
	LogsPath := f.String("logs", "", "evaluate logs from `FILE|API` (e.g. modsec_audit.log or https://loki.example.com)")
	Title := f.String("title", "", "execute only test with `TITLE`")
	Template := f.Bool("template", false, "print tests template and exit")
	Version := f.Bool("version", false, "version")
	RPS := f.Int("rps", 0, "maximum number of requests (tests) per second (e.g. for rate limiting WAFs)")
	Print := f.String("print", "", "print info about tests with status `FAIL|OK|ERR`")

	f.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), usage)
		f.PrintDefaults()
	}

	err := f.Parse(os.Args[1:])
	if err != nil {
		return Flags{}, err
	}

	flags := Flags{
		URL:       stringValue(URL),
		TestsPath: stringValue(TestsPath),
		Verbose:   boolValue(Verbose),
		LogsPath:  stringValue(LogsPath),
		Title:     stringValue(Title),
		Template:  boolValue(Template),
		Version:   boolValue(Version),
		RPS:       intValue(RPS),
		Print:     stringValue(Print),
	}

	err = flags.validate()
	return flags, err
}

func (f Flags) validate() error {
	if !f.Template && !f.Version {
		if f.URL == "" {
			return errors.New("URL cannot be empty")
		}
		if f.TestsPath == "" {
			return errors.New("tests cannot be empty")
		}
	}
	if f.Print != "" &&
		!(f.Print == "FAIL" || f.Print == "OK" || f.Print == "ERR") {
		return errors.New("status must be FAIL, OK or ERR")
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
