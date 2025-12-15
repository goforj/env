//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"os"
)

func main() {
	// IsAppEnvLocal checks if APP_ENV is "local".

	_ = os.Setenv("APP_ENV", env.Local)
	env.Dump(env.IsAppEnvLocal())

	// #bool true
}
