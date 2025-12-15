//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"github.com/goforj/godump"
	"os"
)

func main() {
	// IsAppEnvProduction checks if APP_ENV is "production".

	_ = os.Setenv("APP_ENV", env.Production)
	godump.Dump(env.IsAppEnvProduction())

	// #bool true
}
