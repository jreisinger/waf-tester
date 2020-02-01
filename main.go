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
	help           = flag.Bool("help", false, "print help")
	verbose        = flag.Bool("verbose", false, "be verbose about individual tests")
	all            = flag.Bool("all", false, "print all tests (by default only not OK are printed)")
	only           = flag.String("only", "", "run only these tests (e.g. 920160-1 or ok-tests.txt)")
	testspath      = flag.String("tests", "tests", "directory or file containing tests")
	logspath       = flag.String("logs", "", "file or URL with logs to evaluate (e.g. modsec_audit.log or https://loki.example.com - experimental)")
	stats          = flag.Bool("stats", false, "print progress and statistics about tests")
	tps            = flag.Uint("tps", 10, "tests (HTTP requests) per second")
	scheme         = flag.String("scheme", "https", "http or https scheme")
	waitBeforeEval = flag.Uint("wait", 1, "seconds to wait before evaluating tests")
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
	alltests, err := httptest.GetTests(*testspath, *only, *scheme)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't get tests: %s\n", err)
		os.Exit(1)
	}

	// Check we can get logs.
	if *logspath != "" {
		_, err = httptest.GetLogLines(*logspath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't get logs: %s\n", err)
			os.Exit(1)
		}
	}

	var sepChar = "-"
	var sepLength = 79

	var bar *progressbar.ProgressBar

	var tests []httptest.Test // tests to execute

	for _, test := range alltests {
		// If there are no logs skip tests that don't have exptected status codes.
		if *logspath == "" && len(test.ExpectedStatusCodes) == 0 {
			continue
		}
		tests = append(tests, test)
	}

	if *stats {
		fmt.Printf("LOADED\t%d tests\n", len(alltests))
		fmt.Printf("SKIP\t%d tests\n", len(alltests)-len(tests))
		bar = progressbar.New(len(tests))
	}

	var n uint
	// Execute the tests against the hosts and store results.
	for i := range tests {
		if n >= *tps {
			time.Sleep(1 * time.Second)
			n = 0
		}
		test := &tests[i]
		test.Execute(host)
		if *stats {
			bar.Add(1)
		}
		n++
	}

	if *stats {
		fmt.Println()
		fmt.Println(strings.Repeat(sepChar, sepLength))
	}

	// Logs need to be parsed *after* all requests are done.
	time.Sleep(time.Duration(*waitBeforeEval) * time.Second)
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
