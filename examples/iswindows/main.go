//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"github.com/goforj/godump"
)

func main() {
	// IsWindows reports whether the runtime OS is Windows.

	godump.Dump(env.IsWindows())

	// #bool true  (on Windows)
	// #bool false (elsewhere)
}
