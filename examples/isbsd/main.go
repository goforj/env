//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"github.com/goforj/godump"
)

func main() {
	// IsBSD reports whether the runtime OS is any BSD variant.

	godump.Dump(env.IsBSD())

	// #bool true  (on BSD variants)
	// #bool false (elsewhere)
}
