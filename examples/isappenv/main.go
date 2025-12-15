//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"github.com/goforj/godump"
	"os"
)

func main() {
	// IsAppEnv checks if APP_ENV matches any of the provided environments.

	// Example: match any allowed environment
	_ = os.Setenv("APP_ENV", "staging")
	godump.Dump(env.IsAppEnv(env.Production, env.Staging))

	// #bool true

	// Example: unmatched environment
	_ = os.Setenv("APP_ENV", "local")
	godump.Dump(env.IsAppEnv(env.Production, env.Staging))

	// #bool false
}
