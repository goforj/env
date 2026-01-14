//go:build ignore
// +build ignore

package main

import "github.com/goforj/env/v2"

func main() {
	// IsMac reports whether the runtime OS is macOS (Darwin).

	env.Dump(env.IsMac())

	// #bool true  (on macOS)
	// #bool false (elsewhere)
}
