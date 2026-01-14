//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env/v2"
	"os"
)

func main() {
	// GetUint64 parses a uint64 from an environment variable or fallback string.

	// Example: high range values
	_ = os.Setenv("MAX_ITEMS", "5000")
	maxItems := env.GetUint64("MAX_ITEMS", "100")
	env.Dump(maxItems)

	// #uint64 5000

	// Example: fallback when unset
	os.Unsetenv("MAX_ITEMS")
	maxItems = env.GetUint64("MAX_ITEMS", "100")
	env.Dump(maxItems)

	// #uint64 100
}
