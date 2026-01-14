//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env/v2"
	"os"
)

func main() {
	// GetInt parses an int from an environment variable or fallback string.

	// Example: fallback used
	os.Unsetenv("PORT")
	port := env.GetInt("PORT", "3000")
	env.Dump(port)

	// #int 3000

	// Example: env overrides fallback
	_ = os.Setenv("PORT", "8080")
	port = env.GetInt("PORT", "3000")
	env.Dump(port)

	// #int 8080
}
