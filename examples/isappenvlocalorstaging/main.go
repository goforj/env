//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"github.com/goforj/godump"
	"os"
)

func main() {
	// IsAppEnvLocalOrStaging checks if APP_ENV is either "local" or "staging".

	_ = os.Setenv("APP_ENV", env.Local)
	godump.Dump(env.IsAppEnvLocalOrStaging())

	// #bool true
}
