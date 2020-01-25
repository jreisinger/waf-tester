// Package httptest implements types and functions for testing WAFs.
package httptest

import (
	"bufio"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/jreisinger/waf-tester/yaml"
)

func getOnlyTheseTests(only string) []string {
	var onlyTheseTests []string

	// Is it a file?
	if _, err := os.Stat(only); err == nil {
		f, err := os.Open(only)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			onlyTheseTests = append(onlyTheseTests, line)
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		return onlyTheseTests
	}

	// If not a file it's a test title.
	onlyTheseTests = append(onlyTheseTests, only)
	return onlyTheseTests
}

// GetTests returns the list of available tests.
func GetTests(path string, only string) ([]Test, error) {
	var tests []Test

	// Check path with tests exists.
	if _, err := os.Stat(path); err != nil {
		return tests, err
	}

	var onlyTheseTests []string

	if only != "" {
		onlyTheseTests = getOnlyTheseTests(only)
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

			if len(onlyTheseTests) > 0 && !stringInSlice(t.Title, onlyTheseTests) {
				continue
			}

			t.addCustomHeader()
			tests = append(tests, *t)
		}
	}

	return tests, nil
}

// https://yourbasic.org/golang/generate-uuid-guid/
func genUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

func (t *Test) addCustomHeader() {
	t.ID = genUUID()
	if t.Headers == nil {
		t.Headers = make(map[string]string)
	}
	t.Headers["waf-tester-id"] = t.ID
}

func intInSlice(n int, slice []int) bool {
	for _, m := range slice {
		if n == m {
			return true
		}
	}
	return false
}

func stringInSlice(s string, slice []string) bool {
	for _, r := range slice {
		if s == r {
			return true
		}
	}
	return false
}

// Evaluate sets overall TestStatus to OK|FAIL|ERR. The evaluation is done from:
// * HTTP response status code (primary method)
// * parsing logs
func (t *Test) Evaluate(logspath string) {
	if !t.Executed {
		return
	}

	// There was an error executing the test (HTTP request failed).
	if t.Err != nil {
		if t.ExpectError {
			t.TestStatus = "OK"
		} else {
			t.TestStatus = "ERR"
		}
		return
	}

	if logspath != "" {
		t.AddLogs(logspath)
	}

	// We have output.status defined in the test.
	if len(t.ExpectedStatusCodes) > 0 {
		if intInSlice(t.StatusCode, t.ExpectedStatusCodes) {
			t.TestStatus = "OK"
		} else {
			t.TestStatus = "FAIL"
		}
		return
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

	// No usable output parameters.
	t.Err = errors.New("No output parameters defined in test")
	t.TestStatus = "ERR"
}

// Execute executes a Test. It fills in some of the Test fields (like URL, StatusCode).
func (t *Test) Execute(host string) {
	t.Executed = true

	t.URL = "http" + "://" + path.Join(host, t.Path)

	data := strings.Join(t.Data, "")
	req, err := http.NewRequest(t.Method, t.URL, strings.NewReader(data))
	if err != nil {
		t.Err = err
		return
	}

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
