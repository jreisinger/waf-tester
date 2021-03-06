// Package waftest implements types and functions for testing WAFs.
package waftest

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/jreisinger/waf-tester/wafyaml"
)

// NewHTTPClient creates an HTTP client. Clients and Transports are safe for
// concurrent use by multiple goroutines and for efficiency should only be
// created once and re-used. See https://golang.org/pkg/net/http/ for more.
func NewHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{Timeout: timeout}
}

// GetTests returns the list of available tests.
func GetTests(path string, exec, noexec, header []string, logspath string) ([]Test, error) {
	var tests []Test

	// Check path with tests exists.
	if _, err := os.Stat(path); err != nil {
		return tests, err
	}

	yamls := wafyaml.ParseFiles(path)
	for _, yaml := range yamls {
		for _, test := range yaml.Tests {
			t := &Test{
				Title:               test.Title,
				Desc:                test.Desc,
				Tags:                test.Tags,
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

			// Skip tests that can't be evaluated.
			if logspath == "" && len(t.ExpectedStatusCodes) == 0 {
				continue
			}

			// Skip tests not selected by -exec CLI flag.
			if len(exec) > 0 && (!stringInSlice(t.Title, exec) && !sliceInSlice(exec, t.Tags)) {
				continue
			}

			// Skip tests selected by -noexec CLI flag.
			if len(noexec) > 0 && (stringInSlice(t.Title, noexec) || sliceInSlice(noexec, t.Tags)) {
				continue
			}

			t.addCustomHeaders(header)
			t.setDefaults()

			tests = append(tests, *t)
		}
	}

	return tests, nil
}

func (t *Test) setDefaults() {
	switch {
	case t.Method == "":
		t.Method = "GET"
	case t.Path == "":
		t.Path = "/"
	}
}

// Execute executes a Test. It fills in some of the Test fields (like URL, StatusCode).
func (t *Test) Execute(URL string, client *http.Client) {
	t.Executed = true

	base, err := url.Parse(URL)
	if err != nil {
		t.Err = err
		return
	}
	// fmt.Printf("%#v\n", base)

	u, err := url.Parse(t.Path)
	if err != nil {
		t.Err = err
		return
	}

	t.URL = base.ResolveReference(u).String()
	// fmt.Printf("%#v\n", u)

	data := strings.Join(t.Data, "")
	req, err := http.NewRequest(t.Method, t.URL, strings.NewReader(data))
	if err != nil {
		t.Err = err
		return
	}

	t.fixHostHeader(base.Host)

	for k, v := range t.Headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Err = err
		return
	}

	defer resp.Body.Close()

	// Avoid connections in TIME_WAIT state. See
	// http://tleyden.github.io/blog/2016/11/21/tuning-the-go-http-client-library-for-load-testing/
	// for more.
	_, err = io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		return
	}

	t.StatusCode = resp.StatusCode
	t.Status = resp.Status
}

// Evaluate evaluates a test from response status or logs. Then it sets overall
// TestStatus to OK, FAIL or ERR.
func (t *Test) Evaluate(logspath string) {
	if !t.Executed {
		return
	}
	if t.Err != nil {
		if t.ExpectError { // we were expecting an error
			t.TestStatus = "OK"
		} else {
			t.TestStatus = "ERR"
		}
		return
	}
	if len(t.ExpectedStatusCodes) > 0 {
		t.evaluateFromResponseStatus()
		return
	}
	if t.LogContains != "" || t.LogContainsNot != "" {
		if logspath == "" {
			t.Err = errors.New("can't evaluate test - no logfile supplied (LOG)")
			t.TestStatus = "ERR"
			return
		}
		t.evaluateFromWafLogs()
		return
	}

	t.Err = errors.New("can't evaluate test - no expected (EXP_STATUS, EXP_LOG, EXP_NOLOG) field defined")
	t.TestStatus = "ERR"
}

// evaluateFromResponseStatus evaluates a test by comparing the actual HTTP
// response status code with the expected one.
func (t *Test) evaluateFromResponseStatus() {
	if len(t.ExpectedStatusCodes) == 0 {
		return
	}

	if intInSlice(t.StatusCode, t.ExpectedStatusCodes) {
		t.TestStatus = "OK"
	} else {
		t.TestStatus = "FAIL"
	}
}

// evaluateFromWafLogs evaluates a test by searching expected string in the WAF
// logs.
func (t *Test) evaluateFromWafLogs() {
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
