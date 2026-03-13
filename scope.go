package env

import (
	"os"
	"sort"
	"strings"
	"time"
)

// Scope composes a stable environment variable prefix for related keys.
type Scope struct {
	prefix string
}

// WithPrefix returns a scope rooted at prefix after minimal normalization.
// @group Typed getters
// @behavior readonly
//
// Example: root scope access
//
//	_ = os.Setenv("STORAGE_DRIVER", "local")
//	_ = os.Setenv("STORAGE_ROOT", "storage/app/private")
//
//	storage := env.WithPrefix(" STORAGE ")
//	env.Dump(
//		storage.Key("DRIVER"),
//		storage.Get("DRIVER", "s3"),
//		storage.Get("ROOT", "storage/app/private"),
//	)
//	// #string "STORAGE_DRIVER"
//	// #string "local"
//	// #string "storage/app/private"
func WithPrefix(prefix string) Scope {
	return Scope{prefix: normalizeScopeSegment(prefix)}
}

// Child returns a new scope rooted at the current prefix plus name.
// @group Typed getters
// @behavior readonly
//
// Example: named child scope
//
//	_ = os.Setenv("STORAGE_PUBLIC_ROOT", "storage/app/public")
//
//	public := env.WithPrefix("STORAGE").Child("PUBLIC")
//	env.Dump(
//		public.Key("ROOT"),
//		public.Get("ROOT", "storage/app/public"),
//	)
//	// #string "STORAGE_PUBLIC_ROOT"
//	// #string "storage/app/public"
func (s Scope) Child(name string) Scope {
	child := normalizeScopeSegment(name)
	switch {
	case s.prefix == "":
		return Scope{prefix: child}
	case child == "":
		return s
	default:
		return Scope{prefix: s.prefix + "_" + child}
	}
}

// Key builds the fully qualified environment key for key within the scope.
func (s Scope) Key(key string) string {
	segment := normalizeScopeSegment(key)
	switch {
	case s.prefix == "":
		return segment
	case segment == "":
		return s.prefix
	default:
		return s.prefix + "_" + segment
	}
}

// ChildNames discovers named child scopes under the current prefix.
// @group Typed getters
// @behavior readonly
//
// Example: discover child names
//
//	_ = os.Setenv("STORAGE_DRIVER", "local")
//	_ = os.Setenv("STORAGE_ROOT", "storage/app/private")
//	_ = os.Setenv("STORAGE_PUBLIC_ROOT", "storage/app/public")
//	_ = os.Setenv("STORAGE_AVATARS_BUCKET", "my-bucket")
//	_ = os.Setenv("STORAGE_AVATARS_REGION", "us-east-1")
//
//	names := env.WithPrefix("STORAGE").ChildNames([]string{
//		"DRIVER",
//		"ROOT",
//		"BUCKET",
//		"REGION",
//	})
//	env.Dump(names)
//	// #[]string [
//	//  0 => "AVATARS" #string
//	//  1 => "PUBLIC" #string
//	// ]
func (s Scope) ChildNames(rootKeys []string) []string {
	if s.prefix == "" {
		return []string{}
	}

	rootKeySet := make(map[string]struct{}, len(rootKeys))
	normalizedRootKeys := make([]string, 0, len(rootKeys))
	for _, key := range rootKeys {
		normalized := normalizeScopeSegment(key)
		if normalized == "" {
			continue
		}
		rootKeySet[normalized] = struct{}{}
		normalizedRootKeys = append(normalizedRootKeys, normalized)
	}
	sort.Slice(normalizedRootKeys, func(i, j int) bool {
		if len(normalizedRootKeys[i]) == len(normalizedRootKeys[j]) {
			return normalizedRootKeys[i] < normalizedRootKeys[j]
		}
		return len(normalizedRootKeys[i]) > len(normalizedRootKeys[j])
	})

	prefix := s.prefix + "_"
	children := map[string]struct{}{}

	for _, entry := range os.Environ() {
		key, _, ok := strings.Cut(entry, "=")
		if !ok || !strings.HasPrefix(key, prefix) {
			continue
		}

		remainder := strings.TrimPrefix(key, prefix)
		if remainder == "" {
			continue
		}

		if _, isRootKey := rootKeySet[remainder]; isRootKey {
			continue
		}

		for _, rootKey := range normalizedRootKeys {
			suffix := "_" + rootKey
			if !strings.HasSuffix(remainder, suffix) {
				continue
			}
			child := strings.TrimSuffix(remainder, suffix)
			if child != "" {
				children[child] = struct{}{}
			}
			break
		}
	}

	names := make([]string, 0, len(children))
	for name := range children {
		names = append(names, name)
	}
	sort.Strings(names)

	return names
}

// Get returns the string value for key within the scope.
func (s Scope) Get(key, fallback string) string {
	return Get(s.Key(key), fallback)
}

// GetInt returns the int value for key within the scope.
func (s Scope) GetInt(key, fallback string) int {
	return GetInt(s.Key(key), fallback)
}

// GetInt64 returns the int64 value for key within the scope.
func (s Scope) GetInt64(key, fallback string) int64 {
	return GetInt64(s.Key(key), fallback)
}

// GetUint returns the uint value for key within the scope.
func (s Scope) GetUint(key, fallback string) uint {
	return GetUint(s.Key(key), fallback)
}

// GetUint64 returns the uint64 value for key within the scope.
func (s Scope) GetUint64(key, fallback string) uint64 {
	return GetUint64(s.Key(key), fallback)
}

// GetFloat returns the float64 value for key within the scope.
func (s Scope) GetFloat(key, fallback string) float64 {
	return GetFloat(s.Key(key), fallback)
}

// GetBool returns the bool value for key within the scope.
func (s Scope) GetBool(key, fallback string) bool {
	return GetBool(s.Key(key), fallback)
}

// GetDuration returns the duration value for key within the scope.
func (s Scope) GetDuration(key, fallback string) time.Duration {
	return GetDuration(s.Key(key), fallback)
}

// GetEnum returns the enum value for key within the scope.
func (s Scope) GetEnum(key, fallback string, allowed []string) string {
	return GetEnum(s.Key(key), fallback, allowed)
}

// GetSlice returns the string slice value for key within the scope.
func (s Scope) GetSlice(key, fallback string) []string {
	return GetSlice(s.Key(key), fallback)
}

// GetMap returns the string map value for key within the scope.
func (s Scope) GetMap(key, fallback string) map[string]string {
	return GetMap(s.Key(key), fallback)
}

// GetMapInt returns the int map value for key within the scope.
func (s Scope) GetMapInt(key, fallback string, defaultValue int) map[string]int {
	return GetMapInt(s.Key(key), fallback, defaultValue)
}

func normalizeScopeSegment(segment string) string {
	return strings.Trim(strings.TrimSpace(segment), "_")
}
