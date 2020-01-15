package main

import (
	"fmt"

	"github.com/jreisinger/waf-tester/target"
	"github.com/jreisinger/waf-tester/util"
	"github.com/jreisinger/waf-tester/yaml"
)

type Test struct {
	Desc    string
	Method  string
	URI     string
	Headers map[string]string
}

func main() {
	// Parse YAML files -> URI (Path, Headers)
	yamls := yaml.ParseFiles("yaml")

	var tests []Test

	for _, yaml := range yamls {
		for _, test := range yaml.Tests {
			t := Test{Desc: test.Desc, Method: test.Stages[0].Stage.Input.Method, URI: test.Stages[0].Stage.Input.URI, Headers: test.Stages[0].Stage.Input.Headers}
			tests = append(tests, t)
		}
	}

	ch := make(chan target.Target)
	hosts := util.Flag()

	for _, host := range hosts {
		for _, t := range tests {
			go target.Test(ch, *util.Scheme, host, t.URI, t.Headers)
		}
	}

	format := "%s  (%03.0f): %s\n"

	for range hosts {
		for range tests {
			t := <-ch
			if t.Err != nil {
				fmt.Printf(format, "ERR", float64(t.StatusCode), t.Err)
			} else if t.StatusCode != 403 {
				fmt.Printf(format, "FAIL", float64(t.StatusCode), t.URL)
			} else {
				fmt.Printf(format, "OK", float64(t.StatusCode), t.URL)
			}
		}
	}
}
