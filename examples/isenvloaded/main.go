//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"github.com/goforj/godump"
)

func main() {
	// IsEnvLoaded reports whether LoadEnvFileIfExists was executed in this process.

	godump.Dump(env.IsEnvLoaded())

	// #bool true  (after LoadEnvFileIfExists)
	// #bool false (otherwise)
}
