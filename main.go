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
	log.Printf("flags: %+v", flags)

	// Get tests to execute.
	alltests, err := httptest.GetTests(flags.TestsPath, flags.Only, flags.Scheme)
	if err != nil {
		log.Fatalf("can't get tests: %v", err)
	}
	log.Printf("flags: %+v", alltests)
	/*
		Execute the tests.
		Print the results.
	*/
}
