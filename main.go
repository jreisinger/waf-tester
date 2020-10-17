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

	testsToExecute, err := httptest.GetTests(flags.TestsPath, flags.Title, flags.LogsPath)
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
		for i := range testsToExecute {
			test := &testsToExecute[i]
			testsToExecuteCh <- test
		}
	}()

	testsExecutedCh := make(chan *httptest.Test)

	// Limit the number of requests (tests) per second.
	rate := make(chan bool, flags.RPS)
	for i := 0; i < cap(rate); i++ {
		rate <- true
	}
	// Leaky bucket.
	go func() {
		ticker := time.NewTicker(time.Duration(1000/flags.RPS) * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			_, ok := <-rate
			// If this isn't going to run indefinitely, signal
			// this to return by closing the rate channel.
			if !ok {
				return
			}
		}
	}()

	// Get the tests to execute from the channel. Send the executed ones down
	// another channel. Spawn twice as many workers as the number of tests to
	// execute.
	for i := 0; i < len(testsToExecute)*2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range testsToExecuteCh {
				rate <- true
				t.Execute(flags.URL)
				testsExecutedCh <- t
			}
		}()
	}

	go func() {
		wg.Wait()
		close(testsExecutedCh)
		close(rate)
	}()

	var testsExecuted httptest.Tests

	// Wait for all tests to finish so we can evaluate logs if needed.
	bar := progressbar.Default(int64(len(testsToExecute)), "Running tests")
	for i := range testsExecutedCh {
		bar.Add(1)
		testsExecuted = append(testsExecuted, i)
	}

	if flags.LogsPath != "" {
		testsExecuted.AddLogs(flags.LogsPath)
	}

	// Evaluate and print the tests.
	for _, test := range testsExecuted {
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

	httptest.PrintReport(testsToExecute)
}
