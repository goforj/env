//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env/v2"
	"os"
)

func main() {
	// IsAppEnv checks if APP_ENV matches any of the provided environments.

	// Example: match any allowed environment
	_ = os.Setenv("APP_ENV", "staging")
	env.Dump(env.IsAppEnv(env.Production, env.Staging))

	// #bool true

	// Example: unmatched environment
	_ = os.Setenv("APP_ENV", "local")
	env.Dump(env.IsAppEnv(env.Production, env.Staging))

	// #bool false
}
