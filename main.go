package main

import (
	"flag"
	"fmt"
	"os"

	//"github.com/davecgh/go-spew/spew"
	"github.com/jreisinger/waf-tester/httptest"
)

func init() {
	flag.Usage = func() {
		desc := `Run HTTP tests to evaluate WAF functionality.`
		fmt.Fprintf(os.Stderr, "%s\n\nUsage: %s [options] [host [host2 ...]]\n", desc, os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()
}

var (
	help      = flag.Bool("h", false, "print help")
	verbose   = flag.Bool("v", false, "be verbose")
	all       = flag.Bool("a", false, "print all tests (by default only not OK are printed)")
	only      = flag.String("o", "", "run only these tests (e.g. 920160-1 or ok-tests.txt)")
	testspath = flag.String("t", "tests", "directory or file containing tests")
	logspath  = flag.String("l", "", "file containing WAF logs to evaluate (e.g. modsec_audit.log)")
	stats     = flag.Bool("s", false, "print statistics about tests")
)

func main() {

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	hosts := flag.Args()

	// Execute tests against localhost if no hosts supplied.
	if len(hosts) == 0 {
		hosts = append(hosts, "localhost")
	}

	// Get the tests to execute.
	tests, err := httptest.GetTests(*testspath, *only)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't get tests: %s\n", err)
		os.Exit(1)
	}

	// Execute the tests against the hosts and store results.
	for _, host := range hosts {
		for i := range tests {
			test := &tests[i]
			test.Execute(host)
		}
	}

	// Logs need to be parsed *after* all requests are done.
	for i := range tests {
		test := &tests[i]
		test.Evaluate(*logspath)
	}

	// Print test results.
	for _, test := range tests {
		if *all { // print all tests
			if *verbose {
				test.PrintVerbose()
			} else {
				test.Print()
			}
		} else if test.TestStatus != "OK" { // print only not OK tests
			if *verbose {
				test.PrintVerbose()
				//spew.Dump(test)
			} else {
				test.Print()
			}
		}
	}

	if *stats {
		httptest.PrintStats(tests)
	}
}
