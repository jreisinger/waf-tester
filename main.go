package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jreisinger/waf-tester/httptest"
	"github.com/jreisinger/waf-tester/yaml"
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

	tests, err := httptest.GetTests(flags.TestsPath, flags.Title)
	if err != nil {
		log.Fatalf("cannot get tests: %v", err)
	}

	testsToExecute := make(chan *httptest.Test)

	var wg sync.WaitGroup

	// Send the tests to execute down the channel.
	wg.Add(1)
	go func() {
		for i := range tests {
			test := &tests[i]
			testsToExecute <- test
		}
		close(testsToExecute)
		wg.Done()
	}()

	executedTests := make(chan *httptest.Test)

	// Could be supplied as a flag but we don't want too many flags.
	workers := len(tests) / 100

	if flags.RPS != 0 {
		workers = 1
	}

	// Get the tests to execute from the channel. Send the executed ones down
	// another channel. Start workers goroutines to execute the tests.
	for i := 0; i < workers; i++ {
		wg.Add(1)
		cnt := 0
		go func() {
			for t := range testsToExecute {
				t.Execute(flags.Scheme, flags.Host)
				executedTests <- t
				cnt = cnt + 1
				if flags.RPS != 0 && cnt == flags.RPS {
					time.Sleep(1 * time.Second)
					cnt = 0
				}
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(executedTests)
	}()

	// Evaluate and print the tests.
	for test := range executedTests {
		test.Evaluate(flags.LogsPath)
		if flags.Verbose {
			test.PrintVerbose()
		} else {
			test.Print()
		}
	}

	if flags.Report {
		httptest.PrintReport(tests)
	}
}
