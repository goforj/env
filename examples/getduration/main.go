//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"os"
)

func main() {
	// GetDuration parses a Go duration string (e.g. "5s", "10m", "1h").

	// Example: override request timeout
	_ = os.Setenv("HTTP_TIMEOUT", "30s")
	timeout := env.GetDuration("HTTP_TIMEOUT", "5s")
	env.Dump(timeout)

	// #time.Duration 30s

	// Example: fallback when unset
	os.Unsetenv("HTTP_TIMEOUT")
	timeout = env.GetDuration("HTTP_TIMEOUT", "5s")
	env.Dump(timeout)

	// #time.Duration 5s
}
