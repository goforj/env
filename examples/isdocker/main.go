//go:build ignore
// +build ignore

package main

import "github.com/goforj/env/v2"

func main() {
	// IsDocker reports whether the current process is running in a Docker container.

	// Example: typical host
	env.Dump(env.IsDocker())

	// #bool false (unless inside Docker)
}
