//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"github.com/goforj/godump"
)

func main() {
	// Arch returns the CPU architecture the binary is running on.

	// Example: print GOARCH
	godump.Dump(env.Arch())

	// #string "amd64"
	// #string "arm64"
}
