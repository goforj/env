//go:build ignore
// +build ignore

package main

import "github.com/goforj/env/v2"

func main() {
	// IsLinux reports whether the runtime OS is Linux.

	env.Dump(env.IsLinux())

	// #bool true  (on Linux)
	// #bool false (on other OSes)
}
