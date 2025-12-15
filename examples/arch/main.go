//go:build ignore
// +build ignore

package main

import "github.com/goforj/env"

func main() {
	// Arch returns the CPU architecture the binary is running on.

	// Example: print GOARCH
	env.Dump(env.Arch())

	// #string "amd64"
	// #string "arm64"
}
