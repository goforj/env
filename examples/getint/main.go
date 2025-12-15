//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"github.com/goforj/godump"
	"os"
)

func main() {
	// GetInt parses an int from an environment variable or fallback string.

	// Example: fallback used
	os.Unsetenv("PORT")
	port := env.GetInt("PORT", "3000")
	godump.Dump(port)

	// #int 3000

	// Example: env overrides fallback
	_ = os.Setenv("PORT", "8080")
	port = env.GetInt("PORT", "3000")
	godump.Dump(port)

	// #int 8080
}
