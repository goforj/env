//go:build ignore
// +build ignore

package main

import "github.com/goforj/env"

func main() {
	// IsBSD reports whether the runtime OS is any BSD variant.

	env.Dump(env.IsBSD())

	// #bool true  (on BSD variants)
	// #bool false (elsewhere)
}
