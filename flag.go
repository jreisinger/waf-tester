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
}

func ParseFlags() (Flags, error) {
	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	Host := f.String("host", "", "FQDN or IP address of the host to test")
	TestsPath := f.String("tests", "tests", "directory or file containing tests")
	Scheme := f.String("scheme", "https", "http or https scheme")
	Verbose := f.Bool("verbose", false, "be verbose about individual tests")
	err := f.Parse(os.Args[1:])
	if err != nil {
		return Flags{}, err
	}

	flags := Flags{
		Host:      stringValue(Host),
		TestsPath: stringValue(TestsPath),
		Scheme:    stringValue(Scheme),
		Verbose:   boolValue(Verbose),
	}

	err = flags.validate()
	return flags, err
}

func (f Flags) validate() error {
	if f.Host == "" {
		return errors.New("host cannot be empty")
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
