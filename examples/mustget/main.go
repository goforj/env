//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"os"
)

func main() {
	// MustGet returns the value of key or panics if missing/empty.

	// Example: required secret
	_ = os.Setenv("API_SECRET", "s3cr3t")
	secret := env.MustGet("API_SECRET")
	env.Dump(secret)

	// #string "s3cr3t"

	// Example: panic on missing value
	os.Unsetenv("API_SECRET")
	secret = env.MustGet("API_SECRET") // panics: env variable missing: API_SECRET
}
