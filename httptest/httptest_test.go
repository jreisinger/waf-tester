package httptest

import "testing"

func TestGetTests(t *testing.T) {
	tests, err := GetTests("../tests", "https")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(tests) == 0 {
		t.Errorf("tests: want some tests, got %d tests", len(tests))
	}
}
