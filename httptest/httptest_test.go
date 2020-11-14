package httptest

import "testing"

func TestGetTests(t *testing.T) {
	tests, err := GetTests("../waf_tests", []string{}, []string{}, []string{}, "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(tests) == 0 {
		t.Errorf("tests: want some tests, got %d tests", len(tests))
	}
}
