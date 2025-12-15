//go:build ignore
// +build ignore

package main

import "github.com/goforj/env"

func main() {
	// IsKubernetes reports whether the process is running inside Kubernetes.

	env.Dump(env.IsKubernetes())

	// #bool true  (inside Kubernetes pods)
	// #bool false (elsewhere)
}
