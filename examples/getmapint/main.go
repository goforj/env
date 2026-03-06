//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env/v2"
	"os"
)

func main() {
	// GetMapInt parses key=int pairs separated by commas into a map.
	// Invalid, missing, or non-positive values fall back to defaultValue.

	// Example: parse worker queue weights
	_ = os.Setenv("QUEUE_WEIGHTS", "critical=6, default=3, low=1")
	weights := env.GetMapInt("QUEUE_WEIGHTS", "", 1)
	env.Dump(weights)
	// #map[string]int [
	//  "critical" => 6 #int
	//  "default"  => 3 #int
	//  "low"      => 1 #int
	// ]

	// Example: invalid values use defaultValue
	os.Unsetenv("QUEUE_WEIGHTS")
	weights = env.GetMapInt("QUEUE_WEIGHTS", "critical=,default=0,low=nope,misc", 2)
	env.Dump(weights)
	// #map[string]int [
	//  "critical" => 2 #int
	//  "default"  => 2 #int
	//  "low"      => 2 #int
	//  "misc"     => 2 #int
	// ]
}
