//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"github.com/goforj/godump"
	"os"
)

func main() {
	// GetUint parses a uint from an environment variable or fallback string.

	// Example: defaults to fallback when missing
	os.Unsetenv("WORKERS")
	workers := env.GetUint("WORKERS", "4")
	godump.Dump(workers)

	// #uint 4

	// Example: uses provided unsigned value
	_ = os.Setenv("WORKERS", "16")
	workers = env.GetUint("WORKERS", "4")
	godump.Dump(workers)

	// #uint 16
}
