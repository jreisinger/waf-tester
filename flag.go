package main

import (
	"errors"
	"flag"
	"os"
)

type Flags struct {
	Verbose   bool
	TestsPath string
	LogsPath  string
	Scheme    string
	Host      string
	Title     string
}

func ParseFlags() (Flags, error) {
	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	Host := f.String("host", "", "FQDN or IP address of the host to test")
	TestsPath := f.String("tests", "", "directory or file containing tests")
	Scheme := f.String("scheme", "https", "http or https scheme")
	Verbose := f.Bool("verbose", false, "be verbose about individual tests")
	LogsPath := f.String("logs", "", "filename or API URL with logs to evaluate (modsec_audit.log or https://loki.example.com)")
	Title := f.String("title", "", "execute only test with this title")

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
	}

	err = flags.validate()
	return flags, err
}

func (f Flags) validate() error {
	if f.Host == "" {
		return errors.New("host cannot be empty")
	}
	if f.TestsPath == "" {
		return errors.New("tests cannot be empty")
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
