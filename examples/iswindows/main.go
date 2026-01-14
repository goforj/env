//go:build ignore
// +build ignore

package main

import "github.com/goforj/env/v2"

func main() {
	// IsWindows reports whether the runtime OS is Windows.

	env.Dump(env.IsWindows())

	// #bool true  (on Windows)
	// #bool false (elsewhere)
}
