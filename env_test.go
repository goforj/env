package env

import (
	"os"
	"reflect"
	"testing"
	"time"
)

// Helper: temporarily set an env var and restore after
func withEnv(key, val string, fn func()) {
	original := os.Getenv(key)
	if val == "" {
		_ = os.Unsetenv(key)
	} else {
		_ = os.Setenv(key, val)
	}
	fn()
	_ = os.Setenv(key, original)
}

func TestGet(t *testing.T) {
	withEnv("FOO", "bar", func() {
		if got := Get("FOO", "fallback"); got != "bar" {
			t.Fatalf("expected 'bar', got %q", got)
		}
	})

	withEnv("FOO", "", func() {
		if got := Get("FOO", "fallback"); got != "fallback" {
			t.Fatalf("expected fallback, got %q", got)
		}
	})
}

func TestGetInt(t *testing.T) {
	withEnv("PORT", "8080", func() {
		if got := GetInt("PORT", "1234"); got != 8080 {
			t.Fatalf("expected 8080, got %d", got)
		}
	})
}

func TestGetInt64(t *testing.T) {
	withEnv("MAX", "9223372036854775807", func() {
		if got := GetInt64("MAX", "0"); got != 9223372036854775807 {
			t.Fatalf("int64 mismatch")
		}
	})
}

func TestGetUint(t *testing.T) {
	withEnv("COUNT", "42", func() {
		if got := GetUint("COUNT", "1"); got != 42 {
			t.Fatalf("expected 42, got %d", got)
		}
	})
}

func TestGetUint64(t *testing.T) {
	withEnv("BIGCOUNT", "10000", func() {
		if got := GetUint64("BIGCOUNT", "1"); got != 10000 {
			t.Fatalf("expected 10000, got %d", got)
		}
	})
}

func TestGetFloat(t *testing.T) {
	withEnv("THRESH", "0.75", func() {
		if got := GetFloat("THRESH", "1.0"); got != 0.75 {
			t.Fatalf("expected 0.75, got %f", got)
		}
	})
}

func TestGetBool(t *testing.T) {
	withEnv("DEBUG", "true", func() {
		if !GetBool("DEBUG", "false") {
			t.Fatalf("expected true")
		}
	})
}

func TestGetDuration(t *testing.T) {
	withEnv("TIMEOUT", "5s", func() {
		if got := GetDuration("TIMEOUT", "1s"); got != 5*time.Second {
			t.Fatalf("expected 5s, got %v", got)
		}
	})
}

func TestGetSlice(t *testing.T) {
	withEnv("PEERS", "a,b,c", func() {
		got := GetSlice("PEERS", "")
		expected := []string{"a", "b", "c"}
		if !reflect.DeepEqual(got, expected) {
			t.Fatalf("expected %v, got %v", expected, got)
		}
	})
}

func TestGetMap(t *testing.T) {
	withEnv("LIMITS", "read=10,write=5", func() {
		got := GetMap("LIMITS", "")
		expected := map[string]string{"read": "10", "write": "5"}
		if !reflect.DeepEqual(got, expected) {
			t.Fatalf("expected %v, got %v", expected, got)
		}
	})
}

func TestGetEnum(t *testing.T) {
	withEnv("APP_ENV", "staging", func() {
		got := GetEnum("APP_ENV", "local", []string{"local", "staging", "production"})
		if got != "staging" {
			t.Fatalf("expected staging, got %q", got)
		}
	})
}

func TestGetEnumInvalid(t *testing.T) {
	withEnv("APP_ENV", "invalid", func() {
		if got := GetEnum("APP_ENV", "local", []string{"local", "staging", "production"}); got != "local" {
			t.Fatalf("expected fallback local, got %q", got)
		}
	})
}

func TestGetEnumFallbackNotAllowed(t *testing.T) {
	withEnv("APP_ENV", "unknown", func() {
		if got := GetEnum("APP_ENV", "invalid-fallback", []string{"local", "staging"}); got != "invalid-fallback" {
			t.Fatalf("expected raw fallback when not allowed, got %q", got)
		}
	})
}

func TestSliceAndMapEmptyFallbacks(t *testing.T) {
	withEnv("EMPTY_SLICE", "", func() {
		got := GetSlice("EMPTY_SLICE", "")
		if len(got) != 0 {
			t.Fatalf("expected empty slice, got %v", got)
		}
	})

	withEnv("EMPTY_MAP", "   ", func() {
		got := GetMap("EMPTY_MAP", "")
		if len(got) != 0 {
			t.Fatalf("expected empty map, got %v", got)
		}
	})
}

func TestMustGet(t *testing.T) {
	withEnv("SECRET", "abc123", func() {
		if MustGet("SECRET") != "abc123" {
			t.Fatalf("expected abc123")
		}
	})
}

func TestMustGetMissing(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatalf("expected panic for missing env var")
		}
	}()
	withEnv("MISSING", "", func() {
		MustGet("MISSING")
	})
}

func TestMustGetInt(t *testing.T) {
	withEnv("PORTX", "9000", func() {
		if MustGetInt("PORTX") != 9000 {
			t.Fatalf("expected 9000")
		}
	})
}

func TestMustGetBool(t *testing.T) {
	withEnv("ENABLED", "true", func() {
		if !MustGetBool("ENABLED") {
			t.Fatalf("expected true")
		}
	})
}

func TestGettersReturnFallbackOnBadValues(t *testing.T) {
	withEnv("BAD_INT", "nope", func() {
		if got := GetInt("BAD_INT", "10"); got != 10 {
			t.Fatalf("int fallback expected 10, got %d", got)
		}
		if got := GetInt("BAD_INT", ""); got != 0 {
			t.Fatalf("int fallback expected 0 on bad env+fallback, got %d", got)
		}
	})

	withEnv("BAD_INT64", "xx", func() {
		if got := GetInt64("BAD_INT64", "20"); got != 20 {
			t.Fatalf("int64 fallback expected 20, got %d", got)
		}
		if got := GetInt64("BAD_INT64", ""); got != 0 {
			t.Fatalf("int64 fallback expected 0 on bad env+fallback, got %d", got)
		}
	})

	withEnv("BAD_UINT", "-1", func() {
		if got := GetUint("BAD_UINT", "7"); got != 7 {
			t.Fatalf("uint fallback expected 7, got %d", got)
		}
		if got := GetUint("BAD_UINT", ""); got != 0 {
			t.Fatalf("uint fallback expected 0 on bad env+fallback, got %d", got)
		}
	})

	withEnv("BAD_UINT64", "-1", func() {
		if got := GetUint64("BAD_UINT64", "9"); got != 9 {
			t.Fatalf("uint64 fallback expected 9, got %d", got)
		}
		if got := GetUint64("BAD_UINT64", ""); got != 0 {
			t.Fatalf("uint64 fallback expected 0 on bad env+fallback, got %d", got)
		}
	})

	withEnv("BAD_FLOAT", "not-float", func() {
		if got := GetFloat("BAD_FLOAT", "1.5"); got != 1.5 {
			t.Fatalf("float fallback expected 1.5, got %f", got)
		}
		if got := GetFloat("BAD_FLOAT", ""); got != 0 {
			t.Fatalf("float fallback expected 0 on bad env+fallback, got %f", got)
		}
	})

	withEnv("BAD_BOOL", "maybe", func() {
		if got := GetBool("BAD_BOOL", "true"); got != true {
			t.Fatalf("bool fallback expected true, got %v", got)
		}
		if got := GetBool("BAD_BOOL", ""); got != false {
			t.Fatalf("bool fallback expected false on bad env+fallback, got %v", got)
		}
	})

	withEnv("BAD_DURATION", "sometimes", func() {
		if got := GetDuration("BAD_DURATION", "5s"); got != 5*time.Second {
			t.Fatalf("duration fallback expected 5s, got %v", got)
		}
		if got := GetDuration("BAD_DURATION", ""); got != 0 {
			t.Fatalf("duration fallback expected 0 on bad env+fallback, got %v", got)
		}
	})
}
