//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"github.com/goforj/godump"
	"os"
)

func main() {
	// MustGetBool panics if missing or invalid.

	// Example: gate features explicitly
	_ = os.Setenv("FEATURE_ENABLED", "true")
	enabled := env.MustGetBool("FEATURE_ENABLED")
	godump.Dump(enabled)

	// #bool true

	// Example: panic on invalid value
	_ = os.Setenv("FEATURE_ENABLED", "maybe")
	_ = env.MustGetBool("FEATURE_ENABLED") // panics when parsing
}
