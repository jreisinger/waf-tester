package main

import (
	"os"
	"testing"
)

func TestParseFlags(t *testing.T) {
	rollback := setFlags([]string{"command",
		"-url", "https://example.com",
		"-tests", "tests",
	})
	defer func() { rollback() }()

	flags, err := ParseFlags()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if flags.URL != "https://example.com" {
		t.Errorf("host: want https://example.com, got %s", flags.URL)
	}
	if flags.TestsPath != "tests" {
		t.Errorf("tests: want tests, got %s", flags.TestsPath)
	}
	if flags.Verbose {
		t.Errorf("verbose: want false, got %t", flags.Verbose)
	}
}

// --- helper functions ---

func setFlags(args []string) (rollback func()) {
	osArgs := os.Args
	rollback = func() {
		os.Args = osArgs
	}
	os.Args = args
	return rollback
}
