//go:build ignore
// +build ignore

package main

import "github.com/goforj/env/v2"

func main() {
	// IsEnvLoaded reports whether Load or LoadEnvFileIfExists was executed in this process.

	env.Dump(env.IsEnvLoaded())
	// #bool true  (after Load)
	// #bool false (otherwise)
}
