//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"os"
)

func main() {
	// MustGetInt panics if the value is missing or not an int.

	// Example: ensure numeric port
	_ = os.Setenv("PORT", "8080")
	port := env.MustGetInt("PORT")
	env.Dump(port)

	// #int 8080

	// Example: panic on bad value
	_ = os.Setenv("PORT", "not-a-number")
	_ = env.MustGetInt("PORT") // panics when parsing
}
