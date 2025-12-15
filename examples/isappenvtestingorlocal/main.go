//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"os"
)

func main() {
	// IsAppEnvTestingOrLocal checks if APP_ENV is "testing" or "local".

	_ = os.Setenv("APP_ENV", env.Testing)
	env.Dump(env.IsAppEnvTestingOrLocal())

	// #bool true
}
