//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"github.com/goforj/godump"
	"os"
)

func main() {
	// IsAppEnvStaging checks if APP_ENV is "staging".

	_ = os.Setenv("APP_ENV", env.Staging)
	godump.Dump(env.IsAppEnvStaging())

	// #bool true
}
