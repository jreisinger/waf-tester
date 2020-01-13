package main

import (
	"fmt"
	_ "net/url"

	"github.com/jreisinger/waf-tester/target"
	"github.com/jreisinger/waf-tester/util"
)

var paths = []string{
	"etc/passwd",
	"?page=/etc/passwd",
	"?exec=/bin/bash",
	"?id=1' or '1' = '1'",
	"?<script>",
}

func main() {
	ch := make(chan target.Target)
	hosts := util.Flag()

	for _, host := range hosts {
		for _, path := range paths {
			//path = url.PathEscape(path)
			go target.TargetChecker(ch, *util.Scheme, host, path)
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
