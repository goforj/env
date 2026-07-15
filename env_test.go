package env

import (
	"math/bits"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

// withEnv restores presence as well as value so empty and unset remain distinct between tests.
func withEnv(key, val string, fn func()) {
	original, present := os.LookupEnv(key)
	defer func() {
		if present {
			_ = os.Setenv(key, original)
			return
		}
		_ = os.Unsetenv(key)
	}()
	if val == "" {
		_ = os.Unsetenv(key)
	} else {
		_ = os.Setenv(key, val)
	}
	fn()
}

// TestGet ensures missing strings use the fallback while present values win.
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

// TestGetInt ensures decimal integers parse without changing fallback semantics.
func TestGetInt(t *testing.T) {
	withEnv("PORT", "8080", func() {
		if got := GetInt("PORT", "1234"); got != 8080 {
			t.Fatalf("expected 8080, got %d", got)
		}
	})
}

// TestGetInt64 ensures full-width signed values are preserved.
func TestGetInt64(t *testing.T) {
	withEnv("MAX", "9223372036854775807", func() {
		if got := GetInt64("MAX", "0"); got != 9223372036854775807 {
			t.Fatalf("int64 mismatch")
		}
	})
}

// TestGetUint ensures unsigned values reject invalid text through the fallback path.
func TestGetUint(t *testing.T) {
	withEnv("COUNT", "42", func() {
		if got := GetUint("COUNT", "1"); got != 42 {
			t.Fatalf("expected 42, got %d", got)
		}
	})
}

// TestGetUintUsesNativeWidth ensures overflow behavior follows the target architecture's uint size.
func TestGetUintUsesNativeWidth(t *testing.T) {
	if bits.UintSize != 64 {
		t.Skip("native-width assertion requires a 64-bit uint")
	}
	value := strconv.FormatUint(uint64(1)<<40, 10)
	withEnv("ENV_QPASS_NATIVE_UINT", value, func() {
		if got := GetUint("ENV_QPASS_NATIVE_UINT", "0"); uint64(got) != uint64(1)<<40 {
			t.Fatalf("expected native-width uint, got %d", got)
		}
	})
}

// TestGetUint64 ensures full-width unsigned values are preserved.
func TestGetUint64(t *testing.T) {
	withEnv("BIGCOUNT", "10000", func() {
		if got := GetUint64("BIGCOUNT", "1"); got != 10000 {
			t.Fatalf("expected 10000, got %d", got)
		}
	})
}

// TestGetFloat ensures floating-point values parse without losing fallback behavior.
func TestGetFloat(t *testing.T) {
	withEnv("THRESH", "0.75", func() {
		if got := GetFloat("THRESH", "1.0"); got != 0.75 {
			t.Fatalf("expected 0.75, got %f", got)
		}
	})
}

// TestGetBool ensures standard boolean spellings map predictably.
func TestGetBool(t *testing.T) {
	withEnv("DEBUG", "true", func() {
		if !GetBool("DEBUG", "false") {
			t.Fatalf("expected true")
		}
	})
}

// TestGetDuration ensures Go duration syntax is honored with a safe fallback.
func TestGetDuration(t *testing.T) {
	withEnv("TIMEOUT", "5s", func() {
		if got := GetDuration("TIMEOUT", "1s"); got != 5*time.Second {
			t.Fatalf("expected 5s, got %v", got)
		}
	})
}

// TestGetSlice ensures delimited values are trimmed into stable elements.
func TestGetSlice(t *testing.T) {
	withEnv("PEERS", "a,b,c", func() {
		got := GetSlice("PEERS", "")
		expected := []string{"a", "b", "c"}
		if !reflect.DeepEqual(got, expected) {
			t.Fatalf("expected %v, got %v", expected, got)
		}
	})
}

// TestGetMap ensures delimited key-value text is trimmed and parsed deterministically.
func TestGetMap(t *testing.T) {
	withEnv("LIMITS", " read = 10 ,write= 5, =ignored ", func() {
		got := GetMap("LIMITS", "")
		expected := map[string]string{"read": "10", "write": "5"}
		if !reflect.DeepEqual(got, expected) {
			t.Fatalf("expected %v, got %v", expected, got)
		}
	})
}

// TestGetMapInt ensures numeric map values reject malformed entries through the fallback path.
func TestGetMapInt(t *testing.T) {
	t.Run("parses valid values", func(t *testing.T) {
		withEnv("QUEUE_WEIGHTS", "critical=6, default=3, low=1", func() {
			got := GetMapInt("QUEUE_WEIGHTS", "", 1)
			expected := map[string]int{"critical": 6, "default": 3, "low": 1}
			if !reflect.DeepEqual(got, expected) {
				t.Fatalf("expected %v, got %v", expected, got)
			}
		})
	})

	t.Run("applies default for missing invalid or non-positive values", func(t *testing.T) {
		withEnv("QUEUE_WEIGHTS", "critical=, default=0, low=nope, misc", func() {
			got := GetMapInt("QUEUE_WEIGHTS", "", 2)
			expected := map[string]int{"critical": 2, "default": 2, "low": 2, "misc": 2}
			if !reflect.DeepEqual(got, expected) {
				t.Fatalf("expected %v, got %v", expected, got)
			}
		})
	})

	t.Run("uses fallback string when env is unset", func(t *testing.T) {
		withEnv("QUEUE_WEIGHTS", "", func() {
			got := GetMapInt("QUEUE_WEIGHTS", "critical=9,default=4", 1)
			expected := map[string]int{"critical": 9, "default": 4}
			if !reflect.DeepEqual(got, expected) {
				t.Fatalf("expected %v, got %v", expected, got)
			}
		})
	})

	t.Run("blank input returns empty map", func(t *testing.T) {
		withEnv("QUEUE_WEIGHTS", "   ", func() {
			got := GetMapInt("QUEUE_WEIGHTS", "", 1)
			if len(got) != 0 {
				t.Fatalf("expected empty map, got %v", got)
			}
		})
	})

	t.Run("non-positive default value is normalized to one", func(t *testing.T) {
		withEnv("QUEUE_WEIGHTS", "critical=0,default=nope", func() {
			got := GetMapInt("QUEUE_WEIGHTS", "", 0)
			expected := map[string]int{"critical": 1, "default": 1}
			if !reflect.DeepEqual(got, expected) {
				t.Fatalf("expected %v, got %v", expected, got)
			}
		})
	})
}

// TestGetEnum ensures only explicitly allowed values are returned.
func TestGetEnum(t *testing.T) {
	withEnv("APP_ENV", "staging", func() {
		got := GetEnum("APP_ENV", "local", []string{"local", "staging", "production"})
		if got != "staging" {
			t.Fatalf("expected staging, got %q", got)
		}
	})
}

// TestGetEnumInvalid ensures disallowed configured values fall back safely.
func TestGetEnumInvalid(t *testing.T) {
	withEnv("APP_ENV", "invalid", func() {
		if got := GetEnum("APP_ENV", "local", []string{"local", "staging", "production"}); got != "local" {
			t.Fatalf("expected fallback local, got %q", got)
		}
	})
}

// TestGetEnumFallbackNotAllowed ensures caller fallbacks need not appear in the allowed configured set.
func TestGetEnumFallbackNotAllowed(t *testing.T) {
	withEnv("APP_ENV", "unknown", func() {
		if got := GetEnum("APP_ENV", "invalid-fallback", []string{"local", "staging"}); got != "invalid-fallback" {
			t.Fatalf("expected raw fallback when not allowed, got %q", got)
		}
	})
}

// TestSliceAndMapEmptyFallbacks ensures empty input does not fabricate collection elements.
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

// TestMustGet ensures required strings are returned without fallback ambiguity.
func TestMustGet(t *testing.T) {
	withEnv("SECRET", "abc123", func() {
		if MustGet("SECRET") != "abc123" {
			t.Fatalf("expected abc123")
		}
	})
}

// TestMustGetMissing ensures absent required strings fail fast.
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

// TestMustGetInt ensures required integers return their parsed value.
func TestMustGetInt(t *testing.T) {
	withEnv("PORTX", "9000", func() {
		if MustGetInt("PORTX") != 9000 {
			t.Fatalf("expected 9000")
		}
	})
}

// TestMustGetIntPanicsOnMissingAndInvalid ensures required integer configuration fails fast on absence or corruption.
func TestMustGetIntPanicsOnMissingAndInvalid(t *testing.T) {
	for _, value := range []string{"", "not-an-int"} {
		t.Run(value, func(t *testing.T) {
			withEnv("ENV_QPASS_REQUIRED_INT", value, func() {
				expectPanic(t, "MustGetInt", func() { MustGetInt("ENV_QPASS_REQUIRED_INT") })
			})
		})
	}
}

// TestMustGetBool ensures required booleans return their parsed value.
func TestMustGetBool(t *testing.T) {
	withEnv("ENABLED", "true", func() {
		if !MustGetBool("ENABLED") {
			t.Fatalf("expected true")
		}
	})
}

// TestMustGetBoolPanicsOnMissingAndInvalid ensures required boolean configuration fails fast on absence or corruption.
func TestMustGetBoolPanicsOnMissingAndInvalid(t *testing.T) {
	for _, value := range []string{"", "not-a-bool"} {
		t.Run(value, func(t *testing.T) {
			withEnv("ENV_QPASS_REQUIRED_BOOL", value, func() {
				expectPanic(t, "MustGetBool", func() { MustGetBool("ENV_QPASS_REQUIRED_BOOL") })
			})
		})
	}
}

// FuzzGetMap verifies malformed map entries never panic or produce blank keys.
func FuzzGetMap(f *testing.F) {
	f.Add("read=10, write = 5")
	f.Add("malformed,=empty-key,key=")
	f.Fuzz(func(t *testing.T, value string) {
		for key := range parseStringMap(value) {
			if strings.TrimSpace(key) == "" {
				t.Fatal("GetMap returned an empty key")
			}
		}
	})
}

// TestGettersReturnFallbackOnBadValues ensures optional typed configuration never leaks parse failures as zero values.
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
