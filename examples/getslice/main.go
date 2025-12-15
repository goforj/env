//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"os"
)

func main() {
	// GetSlice splits a comma-separated string into a []string with trimming.

	// Example: trimmed addresses
	_ = os.Setenv("PEERS", "10.0.0.1, 10.0.0.2")
	peers := env.GetSlice("PEERS", "")
	env.Dump(peers)

	// #[]string [
	//  0 => "10.0.0.1" #string
	//  1 => "10.0.0.2" #string
	// ]

	// Example: empty becomes empty slice
	os.Unsetenv("PEERS")
	peers = env.GetSlice("PEERS", "")
	env.Dump(peers)

	// #[]string []
}
