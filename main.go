package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jreisinger/waf-tester/httptest"
	"github.com/jreisinger/waf-tester/yaml"
	"github.com/schollz/progressbar/v3"
)

func init() {
	// Set up CLI tool style of logging.
	log.SetPrefix(os.Args[0] + ": ")
	log.SetFlags(0) // no timestamp
}

// Version is the default version.
var Version = "dev"

func main() {
	flags, err := ParseFlags()
	if err != nil {
		log.Fatalf("cannot parse flags: %v", err)
	}

	if flags.Version {
		fmt.Println(Version)
		os.Exit(0)
	}

	if flags.Template {
		fmt.Print(yaml.Template())
		os.Exit(0)
	}

	tests, err := httptest.GetTests(flags.TestsPath, flags.Title, flags.LogsPath)
	if err != nil {
		log.Fatalf("cannot get tests: %v", err)
	}

	testsToExecuteCh := make(chan *httptest.Test)

	var wg sync.WaitGroup

	// Send the tests to execute down the channel.
	wg.Add(1)
	go func() {
		defer close(testsToExecuteCh)
		defer wg.Done()
		for i := range tests {
			test := &tests[i]
			testsToExecuteCh <- test
		}
	}()

	testsExecutedCh := make(chan *httptest.Test)

	// Could be supplied as a flag but we don't want too many flags.
	workers := len(tests) / 100

	// We want to limit the number of requests (e.g. we are dealing with a rate
	// limiting WAF) or we have less than 100 tests (then workers is 0).
	if flags.RPS != 0 || workers == 0 {
		workers = 1
	}

	// Get the tests to execute from the channel. Send the executed ones down
	// another channel. Start workers goroutines to execute the tests.
	for i := 0; i < workers; i++ {
		wg.Add(1)
		cnt := 0
		go func() {
			defer wg.Done()
			for t := range testsToExecuteCh {
				t.Execute(flags.URL)
				testsExecutedCh <- t
				cnt = cnt + 1
				if flags.RPS != 0 && cnt == flags.RPS {
					time.Sleep(1 * time.Second)
					cnt = 0
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(testsExecutedCh)
	}()

	var doneTests httptest.Tests

	// Wait for all tests to finish so we can evaluate logs if needed.
	bar := progressbar.Default(int64(len(tests)), "Running tests")
	for i := range testsExecutedCh {
		bar.Add(1)
		doneTests = append(doneTests, i)
	}

	if flags.LogsPath != "" {
		doneTests.AddLogs(flags.LogsPath)
	}

	// Evaluate and print the tests.
	for _, test := range doneTests {
		test.Evaluate(flags.LogsPath)

		if !flags.Verbose && flags.Print == "" {
			continue
		}

		if flags.Verbose {
			test.PrintVerbose(flags.Print)
		} else {
			test.Print(flags.Print)
		}
	}

	httptest.PrintReport(tests)
}
