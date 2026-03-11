package env

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestWithPrefixNormalizesPrefix(t *testing.T) {
	scope := WithPrefix(" __STORAGE__ ")
	if got := scope.Key(" __ROOT__ "); got != "STORAGE_ROOT" {
		t.Fatalf("expected STORAGE_ROOT, got %q", got)
	}
}

func TestScopeChildComposition(t *testing.T) {
	scope := WithPrefix("STORAGE").Child(" _PUBLIC_ ")
	if got := scope.Key("ROOT"); got != "STORAGE_PUBLIC_ROOT" {
		t.Fatalf("expected STORAGE_PUBLIC_ROOT, got %q", got)
	}
}

func TestScopeGettersDelegate(t *testing.T) {
	withEnv("STORAGE_DRIVER", "local", func() {
		withEnv("STORAGE_TIMEOUT", "30s", func() {
			withEnv("STORAGE_PUBLIC_ENABLED", "true", func() {
				withEnv("STORAGE_PUBLIC_PEERS", "a,b", func() {
					withEnv("STORAGE_PUBLIC_LIMITS", "read=10,write=5", func() {
						withEnv("STORAGE_PUBLIC_WEIGHTS", "critical=3,default=0", func() {
							storage := WithPrefix("STORAGE")
							public := storage.Child("PUBLIC")

							if got := storage.Get("DRIVER", "s3"); got != "local" {
								t.Fatalf("expected local, got %q", got)
							}
							if got := storage.GetDuration("TIMEOUT", "5s"); got != 30*time.Second {
								t.Fatalf("expected 30s, got %v", got)
							}
							if got := public.GetBool("ENABLED", "false"); !got {
								t.Fatalf("expected true")
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
						})
					})
				})
			})
		})
	})
}

func TestScopeChildNames(t *testing.T) {
	keys := []string{
		"STORAGE_DRIVER",
		"STORAGE_ROOT",
		"STORAGE_PUBLIC_DRIVER",
		"STORAGE_PUBLIC_ROOT",
		"STORAGE_AVATARS_BUCKET",
		"STORAGE_AVATARS_REGION",
		"STORAGE_PUBLIC",
	}

	restore := snapshotEnv(keys)
	defer restore()

	_ = os.Setenv("STORAGE_DRIVER", "local")
	_ = os.Setenv("STORAGE_ROOT", "/tmp/storage")
	_ = os.Setenv("STORAGE_PUBLIC_DRIVER", "local")
	_ = os.Setenv("STORAGE_PUBLIC_ROOT", "/tmp/public")
	_ = os.Setenv("STORAGE_AVATARS_BUCKET", "avatars")
	_ = os.Setenv("STORAGE_AVATARS_REGION", "us-east-1")
	_ = os.Setenv("STORAGE_PUBLIC", "not-a-child")

	names := WithPrefix("STORAGE").ChildNames([]string{
		" DRIVER ",
		"ROOT",
		"BUCKET",
		"REGION",
		"PUBLIC",
	})

	expected := []string{"AVATARS", "PUBLIC"}
	if !reflect.DeepEqual(names, expected) {
		t.Fatalf("expected %v, got %v", expected, names)
	}
}

func TestScopeChildNamesEmptyPrefix(t *testing.T) {
	if got := WithPrefix("___").ChildNames([]string{"ROOT"}); len(got) != 0 {
		t.Fatalf("expected empty names, got %v", got)
	}
}

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
