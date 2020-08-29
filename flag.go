package main

import (
	"errors"
	"flag"
	"os"
)

// Flags are all the available CLI flags (options).
type Flags struct {
	Verbose   bool
	TestsPath string
	LogsPath  string
	Scheme    string
	Host      string
	Title     string
	Report    bool
	Template  bool
	Version   bool
	RPS       int
}

// ParseFlags validates the flags and parses them into Flags.
func ParseFlags() (Flags, error) {
	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	Host := f.String("host", "", "FQDN or IP address of the host to test")
	TestsPath := f.String("tests", "", "directory or file containing tests")
	Scheme := f.String("scheme", "https", "http or https scheme")
	Verbose := f.Bool("verbose", false, "be verbose about individual tests")
	LogsPath := f.String("logs", "", "[EXPERIMENTAL] filename or API URL with logs to evaluate (modsec_audit.log or https://loki.example.com)")
	Title := f.String("title", "", "execute only test with this title")
	Report := f.Bool("report", false, "print overall report about tests")
	Template := f.Bool("template", false, "print tests template and exit")
	Version := f.Bool("version", false, "version")
	RPS := f.Int("rps", 0, "maximum number of requests per second (for rate limiting WAFs)")

	err := f.Parse(os.Args[1:])
	if err != nil {
		return Flags{}, err
	}

	flags := Flags{
		Host:      stringValue(Host),
		TestsPath: stringValue(TestsPath),
		Scheme:    stringValue(Scheme),
		Verbose:   boolValue(Verbose),
		LogsPath:  stringValue(LogsPath),
		Title:     stringValue(Title),
		Report:    boolValue(Report),
		Template:  boolValue(Template),
		Version:   boolValue(Version),
		RPS:       intValue(RPS),
	}

	err = flags.validate()
	return flags, err
}

func (f Flags) validate() error {
	if !f.Template && !f.Version {
		if f.Host == "" {
			return errors.New("host cannot be empty")
		}
		if f.TestsPath == "" {
			return errors.New("tests cannot be empty")
		}
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
