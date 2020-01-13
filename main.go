package main

import (
	"flag"
	"fmt"
	_ "net/url"
	"os"

	"github.com/jreisinger/waf-tester/target"
)

var paths = []string{
	"etc/passwd",
	"?page=/etc/passwd",
	"?exec=/bin/bash",
	"?id=1' or '1' = '1'",
	"?<script>",
}

//
// Command line flags and usage message.
//

var scheme = flag.String("s", "http", "sheme")
var help = flag.Bool("h", false, "print help")

func init() {
	flag.Usage = func() {
		desc := `Test a WAF is blocking malicious requests.`
		fmt.Fprintf(os.Stderr, "%s\n\nUsage: %s [options] host [host2 ...]\n", desc, os.Args[0])
		flag.PrintDefaults()
	}
}

//
// Main.
//

func main() {
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	hosts := flag.Args()
	ch := make(chan target.Target)

	for _, host := range hosts {
		for _, path := range paths {
			//path = url.PathEscape(path)
			go target.TargetChecker(ch, *scheme, host, path)
		}
	}

	for range hosts {
		for range paths {
			t := <-ch
			if t.Err != nil {
				fmt.Printf("ERR  (%03.0f): %s\n", float64(t.StatusCode), t.Err)
			} else if t.StatusCode != 403 {
				fmt.Printf("FAIL (%03.0f): %s\n", float64(t.StatusCode), t.URL)
			} else {
				fmt.Printf("OK   (%03.0f): %s\n", float64(t.StatusCode), t.URL)
			}
		}
	}
}
