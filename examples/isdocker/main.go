//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"github.com/goforj/godump"
)

func main() {
	// IsDocker reports whether the current process is running in a Docker container.

	// Example: typical host
	godump.Dump(env.IsDocker())

	// #bool false (unless inside Docker)
}
