//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"os"
)

func main() {
	// IsAppEnvLocalOrStaging checks if APP_ENV is either "local" or "staging".

	_ = os.Setenv("APP_ENV", env.Local)
	env.Dump(env.IsAppEnvLocalOrStaging())

	// #bool true
}
