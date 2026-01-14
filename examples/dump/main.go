//go:build ignore
// +build ignore

package main

import "github.com/goforj/env/v2"

func main() {
	// Dump is a convenience function that calls godump.Dump.

	// Example: integers
	nums := []int{1, 2, 3}
	env.Dump(nums)

	// #[]int [
	//   0 => 1 #int
	//   1 => 2 #int
	//   2 => 3 #int
	// ]

	// Example: multiple values
	env.Dump("status", map[string]int{"ok": 1, "fail": 0})

	// #string "status"
	// #map[string]int [
	//   "fail" => 0 #int
	//   "ok"   => 1 #int
	// ]
}
