//go:build ignore
// +build ignore

package main

import "github.com/goforj/env"

func main() {
	// IsEnvLoaded reports whether LoadEnvFileIfExists was executed in this process.

	env.Dump(env.IsEnvLoaded())

	// #bool true  (after LoadEnvFileIfExists)
	// #bool false (otherwise)
}
