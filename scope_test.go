package env

import (
	"os"
	"reflect"
	"testing"
	"time"
)

// TestWithPrefixNormalizesPrefix ensures scoped keys use one canonical uppercase separator form.
func TestWithPrefixNormalizesPrefix(t *testing.T) {
	scope := WithPrefix(" __STORAGE__ ")
	if got := scope.Key(" __ROOT__ "); got != "STORAGE_ROOT" {
		t.Fatalf("expected STORAGE_ROOT, got %q", got)
	}
}

// TestScopeChildComposition ensures nested scopes compose without losing ancestor segments.
func TestScopeChildComposition(t *testing.T) {
	scope := WithPrefix("STORAGE").Child(" _PUBLIC_ ")
	if got := scope.Key("ROOT"); got != "STORAGE_PUBLIC_ROOT" {
		t.Fatalf("expected STORAGE_PUBLIC_ROOT, got %q", got)
	}
}

// TestScopeEmptySegmentsPreserveComposition ensures optional empty segments do not introduce malformed separators.
func TestScopeEmptySegmentsPreserveComposition(t *testing.T) {
	if got := WithPrefix("").Child("PUBLIC").Key("ROOT"); got != "PUBLIC_ROOT" {
		t.Fatalf("expected empty root to adopt child, got %q", got)
	}
	if got := WithPrefix("STORAGE").Child("___").Key(""); got != "STORAGE" {
		t.Fatalf("expected empty child and key to preserve root, got %q", got)
	}
	if got := WithPrefix("").Key("ROOT"); got != "ROOT" {
		t.Fatalf("expected unscoped key, got %q", got)
	}
}

// TestScopeGettersDelegate ensures every typed getter resolves through the scoped key.
func TestScopeGettersDelegate(t *testing.T) {
	values := map[string]string{
		"ENV_QPASS_SCOPE_DRIVER":         "local",
		"ENV_QPASS_SCOPE_TIMEOUT":        "30s",
		"ENV_QPASS_SCOPE_INT":            "7",
		"ENV_QPASS_SCOPE_INT64":          "9223372036854775807",
		"ENV_QPASS_SCOPE_UINT":           "8",
		"ENV_QPASS_SCOPE_UINT64":         "18446744073709551615",
		"ENV_QPASS_SCOPE_FLOAT":          "1.5",
		"ENV_QPASS_SCOPE_ENUM":           "blue",
		"ENV_QPASS_SCOPE_PUBLIC_ENABLED": "true",
		"ENV_QPASS_SCOPE_PUBLIC_PEERS":   "a,b",
		"ENV_QPASS_SCOPE_PUBLIC_LIMITS":  "read=10,write=5",
		"ENV_QPASS_SCOPE_PUBLIC_WEIGHTS": "critical=3,default=0",
	}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	restore := snapshotEnv(keys)
	defer restore()
	for key, value := range values {
		_ = os.Setenv(key, value)
	}

	scope := WithPrefix("ENV_QPASS_SCOPE")
	public := scope.Child("PUBLIC")
	if got := scope.Get("DRIVER", "s3"); got != "local" {
		t.Fatalf("expected local, got %q", got)
	}
	if got := scope.GetDuration("TIMEOUT", "5s"); got != 30*time.Second {
		t.Fatalf("expected 30s, got %v", got)
	}
	if got := scope.GetInt("INT", "0"); got != 7 {
		t.Fatalf("expected int 7, got %d", got)
	}
	if got := scope.GetInt64("INT64", "0"); got != 9223372036854775807 {
		t.Fatalf("unexpected int64: %d", got)
	}
	if got := scope.GetUint("UINT", "0"); got != 8 {
		t.Fatalf("expected uint 8, got %d", got)
	}
	if got := scope.GetUint64("UINT64", "0"); got != 18446744073709551615 {
		t.Fatalf("unexpected uint64: %d", got)
	}
	if got := scope.GetFloat("FLOAT", "0"); got != 1.5 {
		t.Fatalf("expected float 1.5, got %v", got)
	}
	if got := scope.GetEnum("ENUM", "red", []string{"red", "blue"}); got != "blue" {
		t.Fatalf("expected enum blue, got %q", got)
	}
	if got := public.GetBool("ENABLED", "false"); !got {
		t.Fatal("expected true")
	}
	if got := public.GetSlice("PEERS", ""); !reflect.DeepEqual(got, []string{"a", "b"}) {
		t.Fatalf("unexpected peers: %v", got)
	}
	if got := public.GetMap("LIMITS", ""); !reflect.DeepEqual(got, map[string]string{"read": "10", "write": "5"}) {
		t.Fatalf("unexpected limits: %v", got)
	}
	if got := public.GetMapInt("WEIGHTS", "", 2); !reflect.DeepEqual(got, map[string]int{"critical": 3, "default": 2}) {
		t.Fatalf("unexpected weights: %v", got)
	}
}

// TestScopeChildNames ensures immediate child discovery is unique and deterministic.
func TestScopeChildNames(t *testing.T) {
	keys := []string{
		"ENV_QPASS_DISCOVERY_DRIVER",
		"ENV_QPASS_DISCOVERY_ROOT",
		"ENV_QPASS_DISCOVERY_PUBLIC_DRIVER",
		"ENV_QPASS_DISCOVERY_PUBLIC_ROOT",
		"ENV_QPASS_DISCOVERY_AVATARS_BUCKET",
		"ENV_QPASS_DISCOVERY_AVATARS_REGION",
		"ENV_QPASS_DISCOVERY_PUBLIC",
	}

	restore := snapshotEnv(keys)
	defer restore()

	_ = os.Setenv("ENV_QPASS_DISCOVERY_DRIVER", "local")
	_ = os.Setenv("ENV_QPASS_DISCOVERY_ROOT", "/tmp/storage")
	_ = os.Setenv("ENV_QPASS_DISCOVERY_PUBLIC_DRIVER", "local")
	_ = os.Setenv("ENV_QPASS_DISCOVERY_PUBLIC_ROOT", "/tmp/public")
	_ = os.Setenv("ENV_QPASS_DISCOVERY_AVATARS_BUCKET", "avatars")
	_ = os.Setenv("ENV_QPASS_DISCOVERY_AVATARS_REGION", "us-east-1")
	_ = os.Setenv("ENV_QPASS_DISCOVERY_PUBLIC", "not-a-child")

	names := WithPrefix("ENV_QPASS_DISCOVERY").ChildNames([]string{
		" DRIVER ",
		"ROOT",
		"__ROOT__",
		"BUCKET",
		"REGION",
		"PUBLIC",
	})

	expected := []string{"AVATARS", "PUBLIC"}
	if !reflect.DeepEqual(names, expected) {
		t.Fatalf("expected %v, got %v", expected, names)
	}
}

// TestScopeChildNamesWithMultiWordChildrenAndRootKeys ensures compound child names survive alongside values on the scope root.
func TestScopeChildNamesWithMultiWordChildrenAndRootKeys(t *testing.T) {
	keys := []string{
		"ENV_QPASS_CACHE_DRIVER",
		"ENV_QPASS_CACHE_PAGE_CACHE_DRIVER",
		"ENV_QPASS_CACHE_PAGE_CACHE_FILE_DIR",
		"ENV_QPASS_CACHE_USER_SESSIONS_DEFAULT_TTL_SECONDS",
		"ENV_QPASS_STORAGE_PUBLIC_S3_ACCESS_KEY_ID",
		"ENV_QPASS_STORAGE_PUBLIC_S3_SECRET_ACCESS_KEY",
	}

	restore := snapshotEnv(keys)
	defer restore()

	_ = os.Setenv("ENV_QPASS_CACHE_DRIVER", "memory")
	_ = os.Setenv("ENV_QPASS_CACHE_PAGE_CACHE_DRIVER", "file")
	_ = os.Setenv("ENV_QPASS_CACHE_PAGE_CACHE_FILE_DIR", "/tmp/page-cache")
	_ = os.Setenv("ENV_QPASS_CACHE_USER_SESSIONS_DEFAULT_TTL_SECONDS", "60")
	_ = os.Setenv("ENV_QPASS_STORAGE_PUBLIC_S3_ACCESS_KEY_ID", "access")
	_ = os.Setenv("ENV_QPASS_STORAGE_PUBLIC_S3_SECRET_ACCESS_KEY", "secret")

	cacheNames := WithPrefix("ENV_QPASS_CACHE").ChildNames([]string{
		"DRIVER",
		"FILE_DIR",
		"DEFAULT_TTL_SECONDS",
	})
	expectedCache := []string{"PAGE_CACHE", "USER_SESSIONS"}
	if !reflect.DeepEqual(cacheNames, expectedCache) {
		t.Fatalf("expected cache child names %v, got %v", expectedCache, cacheNames)
	}

	storageNames := WithPrefix("ENV_QPASS_STORAGE").ChildNames([]string{
		"ACCESS_KEY_ID",
		"SECRET_ACCESS_KEY",
	})
	expectedStorage := []string{"PUBLIC_S3"}
	if !reflect.DeepEqual(storageNames, expectedStorage) {
		t.Fatalf("expected storage child names %v, got %v", expectedStorage, storageNames)
	}
}

// TestScopeChildNamesEmptyPrefix ensures root discovery returns only immediate top-level segments.
func TestScopeChildNamesEmptyPrefix(t *testing.T) {
	if got := WithPrefix("___").ChildNames([]string{"ROOT"}); len(got) != 0 {
		t.Fatalf("expected empty names, got %v", got)
	}
}

// snapshotEnv restores exact process environment presence and values after a test.
func snapshotEnv(keys []string) func() {
	originals := make(map[string]string, len(keys))
	present := make(map[string]bool, len(keys))

	for _, key := range keys {
		value, ok := os.LookupEnv(key)
		if ok {
			originals[key] = value
		}
		present[key] = ok
	}

	return func() {
		for _, key := range keys {
			if present[key] {
				_ = os.Setenv(key, originals[key])
				continue
			}
			_ = os.Unsetenv(key)
		}
	}
}
