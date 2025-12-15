//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"github.com/goforj/godump"
)

func main() {
	// IsDockerInDocker reports whether we are inside a Docker-in-Docker environment.

	godump.Dump(env.IsDockerInDocker())

	// #bool true  (inside DinD containers)
	// #bool false (on hosts or non-DinD containers)
}
