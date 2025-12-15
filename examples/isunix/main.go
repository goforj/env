//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"github.com/goforj/godump"
)

func main() {
	// IsUnix reports whether the OS is Unix-like.

	godump.Dump(env.IsUnix())

	// #bool true  (on Unix-like OSes)
	// #bool false (e.g., on Windows or Plan 9)
}
