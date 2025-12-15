//go:build ignore
// +build ignore

package main

import "github.com/goforj/env"

func main() {
	// IsDockerInDocker reports whether we are inside a Docker-in-Docker environment.

	env.Dump(env.IsDockerInDocker())

	// #bool true  (inside DinD containers)
	// #bool false (on hosts or non-DinD containers)
}
