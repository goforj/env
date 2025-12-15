//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"github.com/goforj/godump"
	"os"
)

func main() {
	// IsAppEnvTesting reports whether APP_ENV is "testing" or the process looks like `go test`.

	// Example: APP_ENV explicitly testing
	_ = os.Setenv("APP_ENV", env.Testing)
	godump.Dump(env.IsAppEnvTesting())

	// #bool true

	// Example: no test markers
	_ = os.Unsetenv("APP_ENV")
	godump.Dump(env.IsAppEnvTesting())

	// #bool false (outside of test binaries)
}
