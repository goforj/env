//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env/v2"
	"os"
)

func main() {
	// Reload re-applies the same layered env loading as Load, even if Load already
	// ran earlier in the same process.

	// Example: refresh changed env files
	_ = os.WriteFile(".env", []byte("SERVICE=api"), 0o644)
	_ = env.Load()
	_ = os.WriteFile(".env", []byte("SERVICE=worker"), 0o644)
	_ = env.Reload()
	env.Dump(os.Getenv("SERVICE"))
	// #string "worker"
}
