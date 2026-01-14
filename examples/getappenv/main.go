//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env/v2"
	"os"
)

func main() {
	// GetAppEnv returns the current APP_ENV (empty string if unset).

	// Example: simple retrieval
	_ = os.Setenv("APP_ENV", "staging")
	env.Dump(env.GetAppEnv())

	// #string "staging"
}
