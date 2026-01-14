//go:build ignore
// +build ignore

package main

import "github.com/goforj/env/v2"

func main() {
	// IsHostEnvironment reports whether the process is running *outside* any
	// container or orchestrated runtime.

	env.Dump(env.IsHostEnvironment())

	// #bool true  (on bare-metal/VM hosts)
	// #bool false (inside containers)
}
