//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env/v2"
	"os"
)

func main() {
	// GetFloat parses a float64 from an environment variable or fallback string.

	// Example: override threshold
	_ = os.Setenv("THRESHOLD", "0.82")
	threshold := env.GetFloat("THRESHOLD", "0.75")
	env.Dump(threshold)

	// #float64 0.82

	// Example: fallback with decimal string
	os.Unsetenv("THRESHOLD")
	threshold = env.GetFloat("THRESHOLD", "0.75")
	env.Dump(threshold)

	// #float64 0.75
}
