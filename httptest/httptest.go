// Package httptest implements types and functions for testing WAFs.
package httptest

import (
	"net/http"
	"time"

	"github.com/jreisinger/waf-tester/yaml"
)

// Test represents an HTTP test. It contains both request and response fields.
type Test struct {
	Desc       string
	Method     string
	Scheme     string
	Host       string
	Path       string // URI
	URL        string // scheme + host + Path
	Headers    map[string]string
	Err        error
	StatusCode int
}

// GetTests returns the list of available tests.
func GetTests(dirname string) []Test {
	var tests []Test

	yamls := yaml.ParseFiles(dirname)
	for _, yaml := range yamls {
		for _, test := range yaml.Tests {
			t := Test{
				Desc:    test.Desc,
				Method:  test.Stages[0].Stage.Input.Method,
				Path:    test.Stages[0].Stage.Input.URI,
				Headers: test.Stages[0].Stage.Input.Headers,
			}
			tests = append(tests, t)
		}
	}

	return tests
}

// Execute executes a Test. It fills in some of the Test fields (like URL, StatusCode).
func (t *Test) Execute(host string) {
	t.URL = "http" + "://" + host + "/" + t.Path

	req, err := http.NewRequest(t.Method, t.URL, nil)
	if err != nil {
		t.Err = err
		return
	}

	for k, v := range t.Headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	if err != nil {
		t.Err = err
		return
	}
	defer resp.Body.Close()

	t.StatusCode = resp.StatusCode
}
