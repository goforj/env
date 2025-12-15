//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"os"
)

func main() {
	// GetBool parses a boolean from an environment variable or fallback string.

	// Example: numeric truthy
	_ = os.Setenv("DEBUG", "1")
	debug := env.GetBool("DEBUG", "false")
	env.Dump(debug)

	// #bool true

	// Example: fallback string
	os.Unsetenv("DEBUG")
	debug = env.GetBool("DEBUG", "false")
	env.Dump(debug)

	// #bool false
}
