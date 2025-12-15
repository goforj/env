//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"github.com/goforj/godump"
	"os"
)

func main() {
	// IsAppEnvLocal checks if APP_ENV is "local".

	_ = os.Setenv("APP_ENV", env.Local)
	godump.Dump(env.IsAppEnvLocal())

	// #bool true
}
