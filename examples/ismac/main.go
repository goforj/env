//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"github.com/goforj/godump"
)

func main() {
	// IsMac reports whether the runtime OS is macOS (Darwin).

	godump.Dump(env.IsMac())

	// #bool true  (on macOS)
	// #bool false (elsewhere)
}
