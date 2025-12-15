//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"github.com/goforj/godump"
	"os"
)

func main() {
	// IsAppEnvDev checks if APP_ENV is "dev".

	_ = os.Setenv("APP_ENV", env.Dev)
	godump.Dump(env.IsAppEnvDev())

	// #bool true
}
