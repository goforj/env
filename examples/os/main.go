//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"github.com/goforj/godump"
)

func main() {
	// OS returns the current operating system identifier.

	// Example: inspect GOOS
	godump.Dump(env.OS())

	// #string "linux"   (on Linux)
	// #string "darwin"  (on macOS)
	// #string "windows" (on Windows)
}
