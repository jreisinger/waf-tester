package util

import (
	"flag"
	"fmt"
	"os"
)

//
// Command line flags and usage message.
//

func init() {
	flag.Usage = func() {
		desc := `Test a WAF is blocking malicious requests.`
		fmt.Fprintf(os.Stderr, "%s\n\nUsage: %s [options] host [host2 ...]\n", desc, os.Args[0])
		flag.PrintDefaults()
	}
}

var Scheme = flag.String("s", "http", "sheme")
var help = flag.Bool("h", false, "print help")

func Flag() []string {
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	return flag.Args()
}
