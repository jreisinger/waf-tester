package main

import (
	"log"
	"os"
	"time"

	"github.com/jreisinger/waf-tester/httptest"
)

func main() {
	// Set up CLI tool style of logging.
	log.SetPrefix(os.Args[0] + ": ")
	log.SetFlags(0) // no timestamp

	flags, err := ParseFlags()
	if err != nil {
		log.Fatalf("cannot parse flags: %v", err)
	}

	tests, err := httptest.GetTests(flags.TestsPath, flags.Scheme, flags.Title)
	if err != nil {
		log.Fatalf("cannot get tests: %v", err)
	}

	for i := range tests {
		test := &tests[i]
		test.Execute(flags.Host)
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
}
