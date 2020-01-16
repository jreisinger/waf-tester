package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jreisinger/waf-tester/httptest"
)

func main() {
	// Set up CLI tool style of logging.
	log.SetPrefix(os.Args[0] + ": ")
	log.SetFlags(0) // no timestamp

	if !(len(os.Args) > 1) {
		// write to stderr and call os.Exit(1)
		log.Fatalf("usage: %s %s", os.Args[0], "host [host2...]")
	}

	hosts := os.Args[1:]

	// Get list of tests to execute.
	tests := httptest.GetTests("tests")
	//spew.Dump(tests)

	// Execute the tests against the hosts and show results.
	for _, host := range hosts {
		for _, test := range tests {
			test.Execute(host)
			//spew.Dump(test)
			fmt.Printf("%03d %-4s %s\n", test.StatusCode, test.Method, test.URL)
		}
	}
}
