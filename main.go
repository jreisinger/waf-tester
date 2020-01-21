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
	help     = flag.Bool("h", false, "print help")
	verbose  = flag.Bool("v", false, "be verbose")
	all      = flag.Bool("a", false, "print all tests (by default only not OK tests are printed)")
	title    = flag.String("t", "", "run only test with this TITLE")
	testsdir = flag.String("d", "tests", "directory containing tests")
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

	tests, err := httptest.GetTests(*testsdir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't get tests: %s\n", err)
		os.Exit(1)
	}

	// Execute the tests against the hosts and show results.
	for _, host := range hosts {
		for _, test := range tests {
			if *title != "" && test.Title != *title {
				continue
			}

			test.Execute(host)

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
}
