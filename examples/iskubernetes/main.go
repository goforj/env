//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"github.com/goforj/godump"
)

func main() {
	// IsKubernetes reports whether the process is running inside Kubernetes.

	godump.Dump(env.IsKubernetes())

	// #bool true  (inside Kubernetes pods)
	// #bool false (elsewhere)
}
