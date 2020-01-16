package main

import (
	"fmt"
	"os"

	"github.com/jreisinger/waf-tester/httptest"
)

func main() {
	// Get list of tests to execute.
	tests := httptest.GetTests("tests")
	//spew.Dump(tests)

	// Execute the tests and show results.
	host := os.Args[1]
	for _, test := range tests {
		test.Execute(host)
		//spew.Dump(test)
		fmt.Printf("%d %s %s\n", test.StatusCode, test.Method, test.Path)
	}

	//	flag := new(util.Flag)
	//	hosts := flag.Parse()
	//
	//	if *flag.ListTests {
	//		for _, t := range tests {
	//			fmt.Printf("%s\n", t.Desc)
	//		}
	//		os.Exit(0)
	//	}
	//
	//	ch := make(chan target.Target)
	//
	//	for _, host := range hosts {
	//		for _, t := range tests {
	//			go target.Test(ch, *flag.Scheme, host, t.URI, t.Headers)
	//		}
	//	}
	//
	//	format := "%s  (%03.0f): %s\n"
	//
	//	for range hosts {
	//		for range tests {
	//			t := <-ch
	//			if t.Err != nil {
	//				fmt.Printf(format, "ERR", float64(t.StatusCode), t.Err)
	//			} else if t.StatusCode != 403 {
	//				fmt.Printf(format, "FAIL", float64(t.StatusCode), t.URL)
	//			} else {
	//				fmt.Printf(format, "OK", float64(t.StatusCode), t.URL)
	//			}
	//		}
	//	}
}
