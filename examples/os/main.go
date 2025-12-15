//go:build ignore
// +build ignore

package main

import "github.com/goforj/env"

func main() {
	// OS returns the current operating system identifier.

	// Example: inspect GOOS
	env.Dump(env.OS())

	// #string "linux"   (on Linux)
	// #string "darwin"  (on macOS)
	// #string "windows" (on Windows)
}
