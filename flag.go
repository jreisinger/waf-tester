package main

import (
	"errors"
	"flag"
	"os"
)

type Flags struct {
	Help           bool
	Verbose        bool
	All            bool
	Only           string
	TestsPath      string
	LogsPath       string
	Stats          bool
	TPS            uint
	Scheme         string
	WaitBeforeEval uint
	Host           string
}

func ParseFlags() (Flags, error) {
	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	Host := f.String("host", "", "FQDN or IP address of the host to test")
	TestsPath := f.String("tests", "tests", "directory or file containing tests")
	Only := f.String("only", "", "run only these tests (e.g. 920160-1 or ok-tests.txt)")
	Scheme := f.String("scheme", "https", "http or https scheme")
	Verbose := f.Bool("verbose", false, "be verbose about individual tests")
	f.Parse(os.Args[1:])

	flags := Flags{
		Host:      stringValue(Host),
		TestsPath: stringValue(TestsPath),
		Only:      stringValue(Only),
		Scheme:    stringValue(Scheme),
		Verbose:   boolValue(Verbose),
	}

	err := flags.validate()
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
