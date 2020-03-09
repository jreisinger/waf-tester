package main

import (
	"log"

	"github.com/jreisinger/waf-tester/httptest"
)

func main() {
	// Get inputs via CLI flags and/or parameters.
	flags, err := ParseFlags()
	if err != nil {
		log.Fatalf("cannot parse flags: %v", err)
	}

	// Get tests to execute.
	alltests, err := httptest.GetTests(flags.TestsPath, flags.Only, flags.Scheme)
	if err != nil {
		log.Fatalf("cannot get tests: %v", err)
	}

	// Execute the tests agains the host.
	for i := range alltests {
		test := &alltests[i]
		test.Execute(flags.Host)
	}

	// Evaluate the tests.
	for i := range alltests {
		test := &alltests[i]
		test.Evaluate(flags.LogsPath)
	}

	// Print the results of the tests.
	for i := range alltests {
		test := alltests[i]
		if flags.Verbose {
			test.PrintVerbose()
		} else {
			test.Print()
		}
	}
}
