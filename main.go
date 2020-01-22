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
		fmt.Fprintf(os.Stderr, "%s\n\nUsage: %s [options] host [host2 ...]\n", desc, os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()
}

var (
	help      = flag.Bool("h", false, "print help")
	verbose   = flag.Bool("v", false, "be verbose")
	all       = flag.Bool("a", false, "print all tests (by default only not OK are printed)")
	only      = flag.String("o", "", "run only test with this TITLE")
	testspath = flag.String("t", "tests", "directory or file containing tests")
	logspath  = flag.String("l", "/tmp/var/log/modsec_audit.log", "file containing WAF logs")
)

func main() {

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	hosts := flag.Args()

	if len(hosts) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	// Get the tests to execute.
	tests, err := httptest.GetTests(*testspath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't get tests: %s\n", err)
		os.Exit(1)
	}

	// Execute the tests against the hosts and store results.
	for _, host := range hosts {
		for i := range tests {
			test := &tests[i]
			if *only != "" && test.Title != *only {
				continue
			}

			test.Execute(host)
		}
	}

	// Logs need to be parsed after all requests are done.
	if *logspath != "" {
		for i := range tests {
			test := &tests[i]
			test.AddLogs(*logspath)
		}
	}

	// Print test results.
	for _, test := range tests {
		if *only != "" && test.Title != *only {
			continue
		}

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
}
