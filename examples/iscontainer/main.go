//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"github.com/goforj/godump"
)

func main() {
	// IsContainer detects common container runtimes (Docker, containerd, Kubernetes, Podman).

	// Example: host vs container
	godump.Dump(env.IsContainer())

	// #bool true  (inside most containers)
	// #bool false (on bare-metal/VM hosts)
}
