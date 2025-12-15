//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"github.com/goforj/godump"
	"os"
)

func main() {
	// IsAppEnvTestingOrLocal checks if APP_ENV is "testing" or "local".

	_ = os.Setenv("APP_ENV", env.Testing)
	godump.Dump(env.IsAppEnvTestingOrLocal())

	// #bool true
}
