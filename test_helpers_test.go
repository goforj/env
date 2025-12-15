package env

import "testing"

// expectPanic asserts that fn panics.
func expectPanic(t *testing.T, name string, fn func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic in %s", name)
		}
	}()
	fn()
}
