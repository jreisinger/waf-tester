// Package httptest implements types and functions for testing WAFs.
package httptest

import (
	"net/http"
	"strings"
	"time"

	"github.com/jreisinger/waf-tester/yaml"
)

// GetTests returns the list of available tests.
func GetTests(dirname string) []Test {
	var tests []Test

	yamls := yaml.ParseFiles(dirname)
	for _, yaml := range yamls {
		for _, test := range yaml.Tests {
			t := &Test{
				Title:   test.Title,
				Desc:    test.Desc,
				File:    test.File,
				Method:  test.Stages[0].Stage.Input.Method,
				Path:    test.Stages[0].Stage.Input.URI,
				Headers: test.Stages[0].Stage.Input.Headers,
				Data:    test.Stages[0].Stage.Input.Data,
			}
			t.setFields()
			tests = append(tests, *t)
		}
	}

	return tests
}

func (t *Test) setFields() {
	switch t.StatusCode {
	case 403:
		t.TestStatus = "OK"
	case 0:
		t.TestStatus = "ERR"
	default:
		t.TestStatus = "FAIL"
	}

	// Set a default value if nil.
	if t.Desc == "" {
		t.Desc = "No test description"
	}
	//if t.Method == "" {
	//	t.Method = "XXX"
	//}
}

// Execute executes a Test. It fills in some of the Test fields (like URL, StatusCode).
func (t *Test) Execute(host string) {
	t.URL = "http" + "://" + host + t.Path

	req, err := http.NewRequest(t.Method, t.URL, strings.NewReader(t.Data))
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
	t.Status = resp.Status
}
