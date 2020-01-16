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
		desc := `Run HTTP tests against hosts protected by a WAF`
		fmt.Fprintf(os.Stderr, "%s\n\nUsage: %s [options] host [host2 ...]\n", desc, os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()
}

var help = flag.Bool("h", false, "print help")
var verbose = flag.Bool("v", false, "be verbose")
var all = flag.Bool("a", false, "print all tests")

func main() {

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	hosts := flag.Args()

	// Get list of tests to execute.
	tests := httptest.GetTests("tests")
	//spew.Dump(tests)

	// Execute the tests against the hosts and show results.
	for _, host := range hosts {
		for _, test := range tests {
			test.Execute(host)

			if *all { // print all tests
				if *verbose {
					test.PrintVerbose()
				} else {
					test.Print()
				}
			} else if test.Status != "OK" { // print only not OK tests
				if *verbose {
					test.PrintVerbose()
				} else {
					test.Print()
				}
			}
		}
	}
}
