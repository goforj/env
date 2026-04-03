//go:build ignore
// +build ignore

package main

import "github.com/goforj/env/v2"

func main() {
	// LoadEnvFileIfExists is a compatibility alias for Load.

	_ = env.LoadEnvFileIfExists()
}
