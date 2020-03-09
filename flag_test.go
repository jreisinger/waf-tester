package main

import (
	"os"
	"testing"
)

func TestParseFlags(t *testing.T) {
	rollback := setFlags([]string{"command",
		"-host", "example.com",
		"-only", "920160-1",
	})
	defer func() { rollback() }()

	flags, err := ParseFlags()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if flags.Host != "example.com" {
		t.Errorf("host: want example.com, got %s", flags.Host)
	}
	if flags.TestsPath != "tests" {
		t.Errorf("tests: want tests, got %s", flags.TestsPath)
	}
	if flags.Only != "920160-1" {
		t.Errorf("only: want 920160-1, got %s", flags.Only)
	}
	if flags.Scheme != "https" {
		t.Errorf("scheme: want https, got %s", flags.Scheme)
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
