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

			format := "%-4s %-9s %s\n"
			if *all { // print all tests
				fmt.Printf(format, test.Status, test.Method, test.URL)
			} else if test.Status != "OK" { // print only not OK tests
				fmt.Printf(format, test.Status, test.Method, test.URL)
			}

			if *verbose {
				//spew.Dump(test)
				fmt.Printf("  %s\n", test.Desc)
				for k, v := range test.Headers {
					fmt.Printf("  %s: %v\n", k, v)
				}
			}
		}
	}
}
