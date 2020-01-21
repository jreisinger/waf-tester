// Package httptest implements types and functions for testing WAFs.
package httptest

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jreisinger/waf-tester/yaml"
)

// GetTests returns the list of available tests.
func GetTests(path string) ([]Test, error) {
	var tests []Test

	// Check directory exists.
	if _, err := os.Stat(path); err != nil {
		return tests, err
	}

	yamls := yaml.ParseFiles(path)
	for _, yaml := range yamls {
		for _, test := range yaml.Tests {
			t := &Test{
				Title:              test.Title,
				Desc:               test.Desc,
				File:               test.File,
				Method:             test.Stages[0].Stage.Input.Method,
				Path:               test.Stages[0].Stage.Input.URI,
				Headers:            test.Stages[0].Stage.Input.Headers,
				Data:               test.Stages[0].Stage.Input.Data,
				ExpectedStatusCode: test.Stages[0].Stage.Output.Status,
			}
			t.setEmptyFields()
			tests = append(tests, *t)
		}
	}

	return tests, nil
}

func (t *Test) setEmptyFields() {
	if t.Desc == "" {
		t.Desc = "No test description"
	}
	if t.ExpectedStatusCode == 0 {
		t.ExpectedStatusCode = 403
	}
	if t.Method == "" {
		t.Method = "XXX"
	}
}

func (t *Test) setTestStatusField() {
	switch t.StatusCode {
	case t.ExpectedStatusCode:
		t.TestStatus = "OK"
	case 0:
		t.TestStatus = "ERR"
	default:
		t.TestStatus = "FAIL"
	}
}

// Execute executes a Test. It fills in some of the Test fields (like URL, StatusCode).
func (t *Test) Execute(host string) {
	t.URL = "http" + "://" + host + t.Path

	defer t.setTestStatusField()

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
