//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"os"
)

func main() {
	// GetUint parses a uint from an environment variable or fallback string.

	// Example: defaults to fallback when missing
	os.Unsetenv("WORKERS")
	workers := env.GetUint("WORKERS", "4")
	env.Dump(workers)

	// #uint 4

	// Example: uses provided unsigned value
	_ = os.Setenv("WORKERS", "16")
	workers = env.GetUint("WORKERS", "4")
	env.Dump(workers)

	// #uint 16
}
