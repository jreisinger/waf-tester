package main

import (
	"log"
	"os"

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

	alltests, err := httptest.GetTests(flags.TestsPath, flags.Scheme)
	if err != nil {
		log.Fatalf("cannot get tests: %v", err)
	}

	for i := range alltests {
		test := &alltests[i]
		test.Execute(flags.Host)
	}

	// Evaluate the tests. This has to be done some time after the tests are
	// executed so it's possible to evaluate the logs.
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
