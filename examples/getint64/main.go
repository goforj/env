//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"os"
)

func main() {
	// GetInt64 parses an int64 from an environment variable or fallback string.

	// Example: parse large numbers safely
	_ = os.Setenv("MAX_SIZE", "1048576")
	size := env.GetInt64("MAX_SIZE", "512")
	env.Dump(size)

	// #int64 1048576

	// Example: fallback when unset
	os.Unsetenv("MAX_SIZE")
	size = env.GetInt64("MAX_SIZE", "512")
	env.Dump(size)

	// #int64 512
}
