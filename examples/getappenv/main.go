//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"github.com/goforj/godump"
	"os"
)

func main() {
	// GetAppEnv returns the current APP_ENV (empty string if unset).

	// Example: simple retrieval
	_ = os.Setenv("APP_ENV", "staging")
	godump.Dump(env.GetAppEnv())

	// #string "staging"
}
