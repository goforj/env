//go:build ignore
// +build ignore

package main

import "github.com/goforj/env/v2"

func main() {
	// IsContainer detects common container runtimes (Docker, containerd, Kubernetes, Podman).

	// Example: host vs container
	env.Dump(env.IsContainer())

	// #bool true  (inside most containers)
	// #bool false (on bare-metal/VM hosts)
}
