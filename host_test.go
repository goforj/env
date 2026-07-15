package env

import (
	"os"
	"testing"
)

// TestIsHostEnvironment_Host ensures absent container evidence identifies a host environment.
func TestIsHostEnvironment_Host(t *testing.T) {
	t.Cleanup(reset)

	// No container indicators
	statFile = func(path string) (os.FileInfo, error) {
		return nil, os.ErrNotExist
	}
	readFile = func(path string) ([]byte, error) {
		return []byte(""), nil
	}
	getEnv = func(s string) string { return "" }

	if !IsHostEnvironment() {
		t.Fatalf("expected IsHostEnvironment == true for host")
	}
}
