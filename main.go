package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jreisinger/waf-tester/waftest"
	"github.com/jreisinger/waf-tester/wafyaml"
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
		log.Fatalf("error parsing flags: %v", err)
	}

	if flags.Version {
		fmt.Println(Version)
		os.Exit(0)
	}

	if flags.Template {
		fmt.Print(wafyaml.Template())
		os.Exit(0)
	}

	testsToExecute, err := waftest.GetTests(flags.TestsPath, flags.Execute, flags.NoExecute, flags.Header, flags.LogsPath)
	if err != nil {
		log.Fatalf("cannot get tests: %v", err)
	}

	testsToExecuteCh := make(chan *waftest.Test)

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

	testsExecutedCh := make(chan *waftest.Test)

	// Limit the number of requests (tests) per second.
	rate := make(chan bool, flags.Rate)
	for i := 0; i < cap(rate); i++ {
		rate <- true
	}
	// Leaky bucket.
	go func() {
		ticker := time.NewTicker(time.Duration(1000/flags.Rate) * time.Millisecond)
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

	// Limit concurrency of requests.
	conc := make(chan bool, flags.Concurrency)

	client := waftest.NewHTTPClient(time.Second * time.Duration(flags.Timeout))

	// Get the tests to execute from the channel. Send the executed ones down
	// another channel. Spawn twice as many workers as the number of tests to
	// execute.
	for i := 0; i < len(testsToExecute)*2; i++ {
		wg.Add(1)
		go func(rate chan bool) {
			defer wg.Done()
			conc <- true
			for t := range testsToExecuteCh {
				rate <- true
				t.Execute(flags.URL, client)
				testsExecutedCh <- t
			}
			<-conc
		}(rate)
	}

	go func() {
		wg.Wait()
		close(testsExecutedCh)
		close(rate)
	}()

	var testsExecuted waftest.Tests

	// Wait for all tests to finish so we can evaluate logs if needed.
	bar := progressbar.Default(int64(len(testsToExecute)), "Executing tests")
	for i := range testsExecutedCh {
		bar.Add(1)
		testsExecuted = append(testsExecuted, i)
	}

	if flags.LogsPath != "" {
		var logsFound, i int
		for {
			i++
			logsFound = testsExecuted.AddLogs(flags.LogsPath)
			if logsFound >= len(testsExecuted) {
				break
			}
			time.Sleep(time.Microsecond * 100) // wait a bit for logs to be written
			if i > 10 {                        // don't wait forever
				break
			}
		}
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

	waftest.PrintReport(testsToExecute)
}
