//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"github.com/goforj/godump"
)

func main() {
	// IsLinux reports whether the runtime OS is Linux.

	godump.Dump(env.IsLinux())

	// #bool true  (on Linux)
	// #bool false (on other OSes)
}
