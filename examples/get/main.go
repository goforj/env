//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env/v2"
	"os"
)

func main() {
	// Get returns the environment variable for key or fallback when empty.

	// Example: fallback when unset
	os.Unsetenv("DB_HOST")
	host := env.Get("DB_HOST", "localhost")
	env.Dump(host)

	// #string "localhost"

	// Example: prefer existing value
	_ = os.Setenv("DB_HOST", "db.internal")
	host = env.Get("DB_HOST", "localhost")
	env.Dump(host)

	// #string "db.internal"
}
