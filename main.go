package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	//"github.com/davecgh/go-spew/spew"
	"github.com/jreisinger/waf-tester/httptest"
	"github.com/schollz/progressbar/v2"
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
	help      = flag.Bool("help", false, "print help")
	verbose   = flag.Bool("verbose", false, "be verbose")
	all       = flag.Bool("all", false, "print all tests (by default only not OK are printed)")
	only      = flag.String("only", "", "run only these tests (e.g. 920160-1 or ok-tests.txt)")
	testspath = flag.String("tests", "tests", "directory or file containing tests")
	logspath  = flag.String("logs", "", "file containing WAF logs to evaluate (e.g. modsec_audit.log)")
	stats     = flag.Bool("stats", false, "print statistics about tests")
	tps       = flag.Uint("tps", 10, "tests (HTTP requests) per second")
	guess     = flag.Bool("guess", false, "try to guess expected test output (result) if not defined in test")
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

	bar := progressbar.New(len(tests))

	var n uint
	// Execute the tests against the hosts and store results.
	for i := range tests {
		if n >= *tps {
			time.Sleep(1 * time.Second)
			n = 0
		}
		test := &tests[i]
		test.Execute(host)
		bar.Add(1)
		n++
	}

	fmt.Println()

	// Logs need to be parsed *after* all requests are done.
	for i := range tests {
		test := &tests[i]
		test.Evaluate(*logspath, *guess)
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
