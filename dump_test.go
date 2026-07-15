package env

import (
	"bytes"
	"os"
	"strings"
	"sync"
	"testing"
)

// TestDumpUsesWriter ensures diagnostics target the caller-provided destination.
func TestDumpUsesWriter(t *testing.T) {
	var buf bytes.Buffer
	setDumpWriter(&buf)
	t.Cleanup(func() { setDumpWriter(os.Stdout) })
	Dump("status", 42)

	out := buf.String()
	if out == "" || !bytes.Contains([]byte(out), []byte("status")) || !bytes.Contains([]byte(out), []byte("42")) {
		t.Fatalf("expected output to include values, got: %q", out)
	}
}

// TestDumpSerializesCompleteWrites ensures concurrent dumps cannot interleave partial environment snapshots.
func TestDumpSerializesCompleteWrites(t *testing.T) {
	var first bytes.Buffer
	var second bytes.Buffer
	setDumpWriter(&first)
	t.Cleanup(func() { setDumpWriter(os.Stdout) })

	const workers = 40
	var wait sync.WaitGroup
	for index := 0; index < workers; index++ {
		wait.Add(1)
		go func(index int) {
			defer wait.Done()
			if index%2 == 0 {
				setDumpWriter(&first)
			} else {
				setDumpWriter(&second)
			}
			Dump("ENV_QPASS_DUMP_MARKER")
		}(index)
	}
	wait.Wait()

	if got := strings.Count(first.String(), "ENV_QPASS_DUMP_MARKER") + strings.Count(second.String(), "ENV_QPASS_DUMP_MARKER"); got != workers {
		t.Fatalf("expected %d complete dumps, got %d", workers, got)
	}
}
