//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"github.com/goforj/godump"
)

func main() {
	// IsDockerHost reports whether this container behaves like a Docker host.

	godump.Dump(env.IsDockerHost())

	// #bool true  (when acting as Docker host)
	// #bool false (for normal containers/hosts)
}
