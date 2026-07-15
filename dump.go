package env

import (
	"io"
	"os"
	"sync"

	"github.com/goforj/godump"
)

var dumpState = struct {
	mu     sync.Mutex
	writer io.Writer
}{writer: os.Stdout}

// setDumpWriter redirects complete dumps so tests can inspect output without partial writes.
func setDumpWriter(w io.Writer) {
	dumpState.mu.Lock()
	defer dumpState.mu.Unlock()
	dumpState.writer = w
}

// Dump writes complete representations of its arguments to standard output.
// @group Debugging
// @behavior readonly
//
// Dump does not redact values. Never pass credentials, tokens, private keys, or other secrets.
//
// Example: integers
//
//	nums := []int{1, 2, 3}
//	env.Dump(nums)
//	// #[]int [
//	//   0 => 1 #int
//	//   1 => 2 #int
//	//   2 => 3 #int
//	// ]
//
// Example: multiple values
//
//	env.Dump("status", map[string]int{"ok": 1, "fail": 0})
//	// #string "status"
//	// #map[string]int [
//	//   "fail" => 0 #int
//	//   "ok"   => 1 #int
//	// ]
func Dump(vs ...any) {
	dumpState.mu.Lock()
	defer dumpState.mu.Unlock()
	d := godump.NewDumper(godump.WithWriter(dumpState.writer))
	d.Dump(vs...)
}
