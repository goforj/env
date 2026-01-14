//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env/v2"
	"os"
)

func main() {
	// IsAppEnvProduction checks if APP_ENV is "production".

	_ = os.Setenv("APP_ENV", env.Production)
	env.Dump(env.IsAppEnvProduction())

	// #bool true
}
