package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/url"
	"os"
	"time"
)

type Target struct {
	Scheme     string
	Host       string
	Path       string
	URL        string
	Err        error
	StatusCode int
}

func (t *Target) getStatusCode() {
	t.URL = t.Scheme + "://" + t.Host + "/" + t.Path

	req, err := http.NewRequest("GET", t.URL, nil)
	if err != nil {
		t.Err = err
		return
	}

	req.Header.Set("Cache-Control", "no-cache")

	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	if err != nil {
		t.Err = err
		return
	}
	defer resp.Body.Close()

	t.StatusCode = resp.StatusCode

	return
}

func TargetChecker(ch chan Target, scheme string, host string, path string) {
	t := Target{Scheme: scheme, Host: host, Path: path}
	t.getStatusCode()
	ch <- t
}

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
	ch := make(chan Target)

	for _, host := range hosts {
		for _, path := range paths {
			//path = url.PathEscape(path)
			go TargetChecker(ch, *scheme, host, path)
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
