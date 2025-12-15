//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"github.com/goforj/godump"
	"os"
)

func main() {
	// GetMap parses key=value pairs separated by commas into a map.

	// Example: parse throttling config
	_ = os.Setenv("LIMITS", "read=10, write=5, burst=20")
	limits := env.GetMap("LIMITS", "")
	godump.Dump(limits)

	// #map[string]string [
	//  "burst" => "20" #string
	//  "read"  => "10" #string
	//  "write" => "5" #string
	// ]

	// Example: returns empty map when unset or blank
	os.Unsetenv("LIMITS")
	limits = env.GetMap("LIMITS", "")
	godump.Dump(limits)

	// #map[string]string []
}
