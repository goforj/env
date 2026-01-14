//go:build ignore
// +build ignore

package main

import "github.com/goforj/env/v2"

func main() {
	// IsContainerOS reports whether this OS is *typically* used as a container base.

	env.Dump(env.IsContainerOS())

	// #bool true  (on Linux)
	// #bool false (on macOS/Windows)
}
