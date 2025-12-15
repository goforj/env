//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"github.com/goforj/godump"
)

func main() {
	// IsHostEnvironment reports whether the process is running *outside* any
	// container or orchestrated runtime.

	godump.Dump(env.IsHostEnvironment())

	// #bool true  (on bare-metal/VM hosts)
	// #bool false (inside containers)
}
