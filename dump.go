package env

import (
	"github.com/goforj/godump"
	"io"
	"os"
)

var dumpWriter io.Writer = os.Stdout

// setDumpWriter allows tests to redirect dump output.
// Not exported — production code never needs this.
func setDumpWriter(w io.Writer) {
	dumpWriter = w
}

// Dump is a convenience function that calls godump.Dump.
// @group Debugging
// @behavior readonly
//
// Example: integers
//
//	nums := []int{1, 2, 3}
//	env.Dump(nums)
//
//	// #[]int [
//	//   0 => 1 #int
//	//   1 => 2 #int
//	//   2 => 3 #int
//	// ]
//
// Example: multiple values
//
//	env.Dump("status", map[string]int{"ok": 1, "fail": 0})
//
//	// #string "status"
//	// #map[string]int [
//	//   "fail" => 0 #int
//	//   "ok"   => 1 #int
//	// ]
func Dump(vs ...any) {
	d := godump.NewDumper(godump.WithWriter(dumpWriter))
	d.Dump(vs...)
}
