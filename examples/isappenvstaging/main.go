//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env/v2"
	"os"
)

func main() {
	// IsAppEnvStaging checks if APP_ENV is "staging".

	_ = os.Setenv("APP_ENV", env.Staging)
	env.Dump(env.IsAppEnvStaging())

	// #bool true
}
