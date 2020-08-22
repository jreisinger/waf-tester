package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jreisinger/waf-tester/httptest"
	"github.com/jreisinger/waf-tester/yaml"
)

func init() {
	// Set up CLI tool style of logging.
	log.SetPrefix(os.Args[0] + ": ")
	log.SetFlags(0) // no timestamp
}

// Version is the default version.
var Version = "dev"

func main() {
	flags, err := ParseFlags()
	if err != nil {
		log.Fatalf("cannot parse flags: %v", err)
	}

	if flags.Version {
		fmt.Println(Version)
		os.Exit(0)
	}

	if flags.Template {
		fmt.Print(yaml.Template())
		os.Exit(0)
	}

	tests, err := httptest.GetTests(flags.TestsPath, flags.Title)
	if err != nil {
		log.Fatalf("cannot get tests: %v", err)
	}

	for i := range tests {
		test := &tests[i]
		test.Execute(flags.Scheme, flags.Host)
	}

	// Let's wait a bit so the logs get written.
	if flags.LogsPath != "" {
		time.Sleep(2 * time.Second)
	}

	// Evaluate the tests. This has to be done some time after the tests are
	// executed so it's possible to evaluate the logs.
	for i := range tests {
		test := &tests[i]
		test.Evaluate(flags.LogsPath)
	}

	// Print the results of the tests.
	for i := range tests {
		test := tests[i]
		if flags.Verbose {
			test.PrintVerbose()
		} else {
			test.Print()
		}
	}

	if flags.Report {
		httptest.PrintReport(tests)
	}
}
