// Package httptest implements types and functions for testing WAFs.
package httptest

import (
	"errors"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/jreisinger/waf-tester/yaml"
)

// GetTests returns the list of available tests.
func GetTests(path string, scheme string) ([]Test, error) {
	var tests []Test

	// Check path with tests exists.
	if _, err := os.Stat(path); err != nil {
		return tests, err
	}

	yamls := yaml.ParseFiles(path)
	for _, yaml := range yamls {
		for _, test := range yaml.Tests {
			t := &Test{
				Title:               test.Title,
				Desc:                test.Desc,
				File:                test.File,
				Method:              test.Stages[0].Stage.Input.Method,
				Path:                test.Stages[0].Stage.Input.URI,
				Headers:             test.Stages[0].Stage.Input.Headers,
				Data:                test.Stages[0].Stage.Input.Data,
				ExpectedStatusCodes: test.Stages[0].Stage.Output.Status,
				LogContains:         test.Stages[0].Stage.Output.LogContains,
				LogContainsNot:      test.Stages[0].Stage.Output.LogContainsNot,
				ExpectError:         test.Stages[0].Stage.Output.ExpectError,
			}

			t.setScheme(scheme)
			t.addCustomHeader()
			tests = append(tests, *t)
		}
	}

	return tests, nil
}

// Execute executes a Test. It fills in some of the Test fields (like URL, StatusCode).
func (t *Test) Execute(host string) {
	t.Executed = true

	t.URL = t.Scheme + "://" + path.Join(host, t.Path)

	data := strings.Join(t.Data, "")
	req, err := http.NewRequest(t.Method, t.URL, strings.NewReader(data))
	if err != nil {
		t.Err = err
		return
	}

	t.fixHostHeader(host)

	for k, v := range t.Headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: time.Second * 2}

	resp, err := client.Do(req)
	if err != nil {
		t.Err = err
		return
	}
	defer resp.Body.Close()

	t.StatusCode = resp.StatusCode
	t.Status = resp.Status
}

// Evaluate sets overall TestStatus to OK|FAIL|ERR.
func (t *Test) Evaluate(logspath string) {
	switch {
	case !t.Executed:
		return
	case t.Err != nil:
		if t.ExpectError { // we were expecting an error
			t.TestStatus = "OK"
		} else {
			t.TestStatus = "ERR"
		}
		return
	case len(t.ExpectedStatusCodes) > 0 || t.LogContains != "" || t.LogContainsNot != "":
		t.EvaluateFromResponseStatus()
		t.EvaluateFromWafLogs(logspath)
	default:
		t.Err = errors.New("no expected (EXP_) fields defined in test")
		t.TestStatus = "ERR"
	}
}

// EvaluateFromResponseStatus evaluates a test by comparing the actual HTTP
// response status code with the expected one.
func (t *Test) EvaluateFromResponseStatus() {
	if intInSlice(t.StatusCode, t.ExpectedStatusCodes) {
		t.TestStatus = "OK"
	} else {
		t.TestStatus = "FAIL"
	}
}

// EvaluateFromWafLogs evaluates a test by searching expected string in the WAF
// logs.
func (t *Test) EvaluateFromWafLogs(logspath string) {
	if logspath != "" {
		t.AddLogs(logspath)
	}

	// We have output.log_contains defined in the test.
	if t.LogContains != "" {
		re := regexp.MustCompile(`\d{6}`) // ex: id "941130"
		id := re.FindString(t.LogContains)
		if foundInLogs(t, id) {
			t.TestStatus = "OK"
		} else {
			t.TestStatus = "FAIL"
		}
		return
	}

	// We have output.no_log_contains defined in the test.
	if t.LogContainsNot != "" {
		re := regexp.MustCompile(`\d{6}`) // ex: id "941130"
		id := re.FindString(t.LogContainsNot)
		if !foundInLogs(t, id) {
			t.TestStatus = "OK"
		} else {
			t.TestStatus = "FAIL"
		}
		return
	}
}
