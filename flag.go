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
	f.Parse(os.Args[1:])

	flags := Flags{
		Host: stringValue(Host),
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
