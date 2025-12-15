package env

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Get returns the environment variable for key or fallback when empty.
// @group Typed getters
// @behavior readonly
//
// Examples use github.com/goforj/godump to illustrate the concrete type.
//
// Example: fallback when unset
//
//	os.Unsetenv("DB_HOST")
//	host := env.Get("DB_HOST", "localhost")
//	env.Dump(host)
//
//	// #string "localhost"
//
// Example: prefer existing value
//
//	_ = os.Setenv("DB_HOST", "db.internal")
//	host = env.Get("DB_HOST", "localhost")
//	env.Dump(host)
//
//	// #string "db.internal"
func Get(key, fallback string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		return fallback
	}
	return val
}

// GetInt parses an int from an environment variable or fallback string.
// @group Typed getters
// @behavior panic
//
// Panics if the chosen value cannot be parsed as base-10 int.
//
// Example: fallback used
//
//	os.Unsetenv("PORT")
//	port := env.GetInt("PORT", "3000")
//	env.Dump(port)
//
//	// #int 3000
//
// Example: env overrides fallback
//
//	_ = os.Setenv("PORT", "8080")
//	port = env.GetInt("PORT", "3000")
//	env.Dump(port)
//
//	// #int 8080
func GetInt(key, fallback string) int {
	val := Get(key, fallback)
	ret, err := strconv.Atoi(val)
	if err != nil {
		panic(err)
	}
	return ret
}

// GetInt64 parses an int64 from an environment variable or fallback string.
// @group Typed getters
// @behavior panic
//
// Example: parse large numbers safely
//
//	_ = os.Setenv("MAX_SIZE", "1048576")
//	size := env.GetInt64("MAX_SIZE", "512")
//	env.Dump(size)
//
//	// #int64 1048576
//
// Example: fallback when unset
//
//	os.Unsetenv("MAX_SIZE")
//	size = env.GetInt64("MAX_SIZE", "512")
//	env.Dump(size)
//
//	// #int64 512
func GetInt64(key, fallback string) int64 {
	val := Get(key, fallback)
	ret, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		panic(err)
	}
	return ret
}

// GetUint parses a uint from an environment variable or fallback string.
// @group Typed getters
// @behavior panic
//
// Example: defaults to fallback when missing
//
//	os.Unsetenv("WORKERS")
//	workers := env.GetUint("WORKERS", "4")
//	env.Dump(workers)
//
//	// #uint 4
//
// Example: uses provided unsigned value
//
//	_ = os.Setenv("WORKERS", "16")
//	workers = env.GetUint("WORKERS", "4")
//	env.Dump(workers)
//
//	// #uint 16
func GetUint(key, fallback string) uint {
	val := Get(key, fallback)
	i, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		panic(err)
	}
	return uint(i)
}

// GetUint64 parses a uint64 from an environment variable or fallback string.
// @group Typed getters
// @behavior panic
//
// Example: high range values
//
//	_ = os.Setenv("MAX_ITEMS", "5000")
//	maxItems := env.GetUint64("MAX_ITEMS", "100")
//	env.Dump(maxItems)
//
//	// #uint64 5000
//
// Example: fallback when unset
//
//	os.Unsetenv("MAX_ITEMS")
//	maxItems = env.GetUint64("MAX_ITEMS", "100")
//	env.Dump(maxItems)
//
//	// #uint64 100
func GetUint64(key, fallback string) uint64 {
	val := Get(key, fallback)
	i, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

// GetFloat parses a float64 from an environment variable or fallback string.
// @group Typed getters
// @behavior panic
//
// Example: override threshold
//
//	_ = os.Setenv("THRESHOLD", "0.82")
//	threshold := env.GetFloat("THRESHOLD", "0.75")
//	env.Dump(threshold)
//
//	// #float64 0.82
//
// Example: fallback with decimal string
//
//	os.Unsetenv("THRESHOLD")
//	threshold = env.GetFloat("THRESHOLD", "0.75")
//	env.Dump(threshold)
//
//	// #float64 0.75
func GetFloat(key, fallback string) float64 {
	val := Get(key, fallback)
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		panic(err)
	}
	return f
}

// GetBool parses a boolean from an environment variable or fallback string.
// @group Typed getters
// @behavior panic
//
// Accepted values: true/false, 1/0, t/f (case-insensitive).
//
// Example: numeric truthy
//
//	_ = os.Setenv("DEBUG", "1")
//	debug := env.GetBool("DEBUG", "false")
//	env.Dump(debug)
//
//	// #bool true
//
// Example: fallback string
//
//	os.Unsetenv("DEBUG")
//	debug = env.GetBool("DEBUG", "false")
//	env.Dump(debug)
//
//	// #bool false
func GetBool(key, fallback string) bool {
	val := Get(key, fallback)
	ret, err := strconv.ParseBool(val)
	if err != nil {
		panic(err)
	}
	return ret
}

// GetDuration parses a Go duration string (e.g. "5s", "10m", "1h").
// @group Typed getters
// @behavior panic
//
// Example: override request timeout
//
//	_ = os.Setenv("HTTP_TIMEOUT", "30s")
//	timeout := env.GetDuration("HTTP_TIMEOUT", "5s")
//	env.Dump(timeout)
//
//	// #time.Duration 30s
//
// Example: fallback when unset
//
//	os.Unsetenv("HTTP_TIMEOUT")
//	timeout = env.GetDuration("HTTP_TIMEOUT", "5s")
//	env.Dump(timeout)
//
//	// #time.Duration 5s
func GetDuration(key, fallback string) time.Duration {
	val := Get(key, fallback)
	d, err := time.ParseDuration(val)
	if err != nil {
		panic(err)
	}
	return d
}

// GetSlice splits a comma-separated string into a []string with trimming.
// @group Typed getters
// @behavior readonly
//
// Example: trimmed addresses
//
//	_ = os.Setenv("PEERS", "10.0.0.1, 10.0.0.2")
//	peers := env.GetSlice("PEERS", "")
//	env.Dump(peers)
//
//	// #[]string [
//	//  0 => "10.0.0.1" #string
//	//  1 => "10.0.0.2" #string
//	// ]
//
// Example: empty becomes empty slice
//
//	os.Unsetenv("PEERS")
//	peers = env.GetSlice("PEERS", "")
//	env.Dump(peers)
//
//	// #[]string []
func GetSlice(key, fallback string) []string {
	val := Get(key, fallback)
	if val == "" {
		return []string{}
	}

	parts := strings.Split(val, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

// GetMap parses key=value pairs separated by commas into a map.
// @group Typed getters
// @behavior readonly
//
// Example: parse throttling config
//
//	_ = os.Setenv("LIMITS", "read=10, write=5, burst=20")
//	limits := env.GetMap("LIMITS", "")
//	env.Dump(limits)
//
//	// #map[string]string [
//	//  "burst" => "20" #string
//	//  "read"  => "10" #string
//	//  "write" => "5" #string
//	// ]
//
// Example: returns empty map when unset or blank
//
//	os.Unsetenv("LIMITS")
//	limits = env.GetMap("LIMITS", "")
//	env.Dump(limits)
//
//	// #map[string]string []
func GetMap(key, fallback string) map[string]string {
	val := Get(key, fallback)
	m := map[string]string{}

	if strings.TrimSpace(val) == "" {
		return m
	}

	pairs := strings.Split(val, ",")
	for _, p := range pairs {
		kv := strings.SplitN(strings.TrimSpace(p), "=", 2)
		if len(kv) == 2 {
			m[kv[0]] = kv[1]
		}
	}

	return m
}

// GetEnum ensures the environment variable's value is in the allowed list.
// @group Typed getters
// @behavior panic
//
// Panic occurs when the chosen value is not in the allowed slice (case-sensitive).
//
// Example: accept only staged environments
//
//	_ = os.Setenv("APP_ENV", "prod")
//	appEnv := env.GetEnum("APP_ENV", "dev", []string{"dev", "staging", "prod"})
//	env.Dump(appEnv)
//
//	// #string "prod"
//
// Example: fallback when unset
//
//	os.Unsetenv("APP_ENV")
//	appEnv = env.GetEnum("APP_ENV", "dev", []string{"dev", "staging", "prod"})
//	env.Dump(appEnv)
//
//	// #string "dev"
func GetEnum(key, fallback string, allowed []string) string {
	val := Get(key, fallback)
	for _, a := range allowed {
		if val == a {
			return val
		}
	}
	panic("env: invalid enum value for " + key + ": " + val)
}

// MustGet returns the value of key or panics if missing/empty.
// @group Typed getters
// @behavior panic
//
// Example: required secret
//
//	_ = os.Setenv("API_SECRET", "s3cr3t")
//	secret := env.MustGet("API_SECRET")
//	env.Dump(secret)
//
//	// #string "s3cr3t"
//
// Example: panic on missing value
//
//	os.Unsetenv("API_SECRET")
//	secret = env.MustGet("API_SECRET") // panics: env variable missing: API_SECRET
func MustGet(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic("env variable missing: " + key)
	}
	return val
}

// MustGetInt panics if the value is missing or not an int.
// @group Typed getters
// @behavior panic
//
// Example: ensure numeric port
//
//	_ = os.Setenv("PORT", "8080")
//	port := env.MustGetInt("PORT")
//	env.Dump(port)
//
//	// #int 8080
//
// Example: panic on bad value
//
//	_ = os.Setenv("PORT", "not-a-number")
//	_ = env.MustGetInt("PORT") // panics when parsing
func MustGetInt(key string) int {
	return GetInt(key, "")
}

// MustGetBool panics if missing or invalid.
// @group Typed getters
// @behavior panic
//
// Example: gate features explicitly
//
//	_ = os.Setenv("FEATURE_ENABLED", "true")
//	enabled := env.MustGetBool("FEATURE_ENABLED")
//	env.Dump(enabled)
//
//	// #bool true
//
// Example: panic on invalid value
//
//	_ = os.Setenv("FEATURE_ENABLED", "maybe")
//	_ = env.MustGetBool("FEATURE_ENABLED") // panics when parsing
func MustGetBool(key string) bool {
	return GetBool(key, "")
}
