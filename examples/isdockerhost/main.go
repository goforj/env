//go:build ignore
// +build ignore

package main

import "github.com/goforj/env"

func main() {
	// IsDockerHost reports whether this container behaves like a Docker host.

	env.Dump(env.IsDockerHost())

	// #bool true  (when acting as Docker host)
	// #bool false (for normal containers/hosts)
}
