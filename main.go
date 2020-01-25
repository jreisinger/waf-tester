package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	//"github.com/davecgh/go-spew/spew"
	"github.com/jreisinger/waf-tester/httptest"
)

func init() {
	flag.Usage = func() {
		desc := `Run HTTP tests to evaluate WAF functionality.`
		fmt.Fprintf(os.Stderr, "%s\n\nUsage: %s [options] [host]\n", desc, os.Args[0])
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
	tps       = flag.Uint("tps", 10, "tests (HTTP requests) per second")
)

func main() {

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	host := "localhost"

	// Execute tests against localhost if no hosts supplied.
	if len(flag.Args()) == 1 {
		host = flag.Arg(0)
	} else if len(flag.Args()) > 1 {
		flag.Usage()
		os.Exit(1)
	}

	// Get the tests to execute.
	tests, err := httptest.GetTests(*testspath, *only)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't get tests: %s\n", err)
		os.Exit(1)
	}

	var sepChar = "-"
	var sepLength = 79

	if *verbose {
		fmt.Println(strings.Repeat(sepChar, sepLength))
		fmt.Printf("Running %d tests against %s\n", len(tests), host)
		fmt.Println(strings.Repeat(sepChar, sepLength))
	}

	throttle := time.Tick(time.Second) // stop for a second

	var n uint
	// Execute the tests against the hosts and store results.
	for i := range tests {
		if n >= *tps {
			<-throttle
			n = 0
		}
		test := &tests[i]
		test.Execute(host)
		n++
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
		fmt.Println(strings.Repeat(sepChar, sepLength))
		httptest.PrintStats(tests)
	}
}
