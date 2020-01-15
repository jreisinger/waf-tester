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

type Flag struct {
	Scheme    *string
	ListTests *bool
}

func (f *Flag) Parse() []string {
	f.Scheme = flag.String("s", "http", "sheme")
	f.ListTests = flag.Bool("l", false, "list tests")

	var help = flag.Bool("h", false, "print help")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	return flag.Args()
}
