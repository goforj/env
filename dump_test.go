package env

import (
	"bytes"
	"testing"
)

func TestDumpUsesWriter(t *testing.T) {
	original := dumpWriter
	defer func() { dumpWriter = original }()

	var buf bytes.Buffer
	setDumpWriter(&buf)
	Dump("status", 42)

	out := buf.String()
	if out == "" || !bytes.Contains([]byte(out), []byte("status")) || !bytes.Contains([]byte(out), []byte("42")) {
		t.Fatalf("expected output to include values, got: %q", out)
	}
}
