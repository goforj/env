//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"github.com/goforj/godump"
)

func main() {
	// IsContainerOS reports whether this OS is *typically* used as a container base.

	godump.Dump(env.IsContainerOS())

	// #bool true  (on Linux)
	// #bool false (on macOS/Windows)
}
