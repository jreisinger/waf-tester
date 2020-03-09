package main

import (
	"os"
	"testing"
)

func TestParseFlags(t *testing.T) {
	rollback := setFlags([]string{
		"command", "--host", "example.com",
	})
	defer func() { rollback() }()

	flags, err := ParseFlags()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if flags.Host != "example.com" {
		t.Errorf("host: want example.com, got %s", flags.Host)
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
