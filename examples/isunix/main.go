//go:build ignore
// +build ignore

package main

import "github.com/goforj/env/v2"

func main() {
	// IsUnix reports whether the OS is Unix-like.

	env.Dump(env.IsUnix())

	// #bool true  (on Unix-like OSes)
	// #bool false (e.g., on Windows or Plan 9)
}
