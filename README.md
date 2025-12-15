<p align="center">
  <img src="./docs/images/logo.png" width="600" alt="goforj/env logo">
</p>

<p align="center">
    Typed environment variables for Go – safe defaults, app env helpers, and zero-ceremony configuration.
</p>

<p align="center">
    <a href="https://pkg.go.dev/github.com/goforj/env"><img src="https://pkg.go.dev/badge/github.com/goforj/env.svg" alt="Go Reference"></a>
    <a href="LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License: MIT"></a>
    <a href="https://github.com/goforj/env/actions"><img src="https://github.com/goforj/env/actions/workflows/test.yml/badge.svg" alt="Go Test"></a>
    <a href="https://golang.org"><img src="https://img.shields.io/badge/go-1.18+-blue?logo=go" alt="Go version"></a>
    <img src="https://img.shields.io/github/v/tag/goforj/env?label=version&sort=semver" alt="Latest tag">
    <a href="https://codecov.io/gh/goforj/env" ><img src="https://codecov.io/github/goforj/env/graph/badge.svg?token=M7EFUVV1XW"/></a>
    <a href="https://goreportcard.com/report/github.com/goforj/env"><img src="https://goreportcard.com/badge/github.com/goforj/env" alt="Go Report Card"></a>
</p>

<p align="center">
  <code>env</code> provides strongly-typed access to environment variables with predictable fallbacks.  
  Eliminate string parsing, centralize app environment checks, and keep configuration boring.  
  Designed to feel native to Go - and invisible when things are working.
</p>

# Features

- 🔐 **Strongly typed getters** – `int`, `bool`, `float`, `duration`, slices, maps
- 🧯 **Safe fallbacks** – never panic, never accidentally empty
- 🌎 **Application environment helpers** – `dev`, `local`, `prod`
- 🧩 **Zero dependencies** – pure Go, lightweight
- 🧭 **Framework-agnostic** – works with any Go app
- 📐 **Enum validation** – constrain values with allowed sets
- 🧼 **Predictable behavior** – no magic, no global state surprises
- 🧱 **Composable building block** – ideal for config structs and startup wiring

## Why env?

Accessing environment variables in Go often leads to:

- Repeated parsing logic
- Unsafe string conversions
- Inconsistent defaults
- Scattered app environment checks

**env** solves this by providing **typed accessors with fallbacks**, so configuration stays boring and predictable.

## Features

- Strongly typed getters (`int`, `bool`, `duration`, slices, maps)
- Safe fallbacks (never panic, never empty by accident)
- App environment helpers (`dev`, `local`, `prod`)
- Zero dependencies
- Framework-agnostic

## Installation

```bash
go get github.com/goforj/env
````

## Usage

```go
import "github.com/goforj/env"

port := env.GetInt("PORT", 8080)
debug := env.GetBool("DEBUG", false)
timeout := env.GetDuration("REQUEST_TIMEOUT", time.Second*5)
```

### Environment file loading

This package uses `github.com/joho/godotenv` for `.env` file loading.
It is intentionally composed into the runtime detection and APP_ENV model rather than reimplemented.

### Application environment

```go
if env.IsAppEnvDev() {
    // dev-only behavior
}
```

<!-- api:embed:start -->

### Index

| Group | Functions |
|------:|-----------|
| **Application environment** | [GetAppEnv](#getappenv) [IsAppEnv](#isappenv) [IsAppEnvDev](#isappenvdev) [IsAppEnvLocal](#isappenvlocal) [IsAppEnvLocalOrStaging](#isappenvlocalorstaging) [IsAppEnvProduction](#isappenvproduction) [IsAppEnvStaging](#isappenvstaging) [IsAppEnvTesting](#isappenvtesting) [IsAppEnvTestingOrLocal](#isappenvtestingorlocal) |
| **Other** | [IsContainer](#iscontainer) [IsDocker](#isdocker) [IsDockerHost](#isdockerhost) [IsDockerInDocker](#isdockerindocker) [IsEnvLoaded](#isenvloaded) [IsHostEnvironment](#ishostenvironment) [IsKubernetes](#iskubernetes) [LoadEnvFileIfExists](#loadenvfileifexists) |
| **Runtime** | [Arch](#arch) [IsBSD](#isbsd) [IsContainerOS](#iscontaineros) [IsLinux](#islinux) [IsMac](#ismac) [IsUnix](#isunix) [IsWindows](#iswindows) [OS](#os) |
| **Typed getters** | [Get](#get) [GetBool](#getbool) [GetDuration](#getduration) [GetEnum](#getenum) [GetFloat](#getfloat) [GetInt](#getint) [GetInt64](#getint64) [GetMap](#getmap) [GetSlice](#getslice) [GetUint](#getuint) [GetUint64](#getuint64) [MustGet](#mustget) [MustGetBool](#mustgetbool) [MustGetInt](#mustgetint) |


## Application environment

### <a id="getappenv"></a>GetAppEnv · readonly

GetAppEnv returns the current APP_ENV (empty string if unset).

_Example: simple retrieval_

```go
_ = os.Setenv("APP_ENV", "staging")
godump.Println(env.GetAppEnv())

// #string "staging"
```

### <a id="isappenv"></a>IsAppEnv · readonly

IsAppEnv checks if APP_ENV matches any of the provided environments.

_Example: match any allowed environment_

```go
_ = os.Setenv("APP_ENV", "staging")
godump.Println(env.IsAppEnv(env.Production, env.Staging))

// #bool true
```

_Example: unmatched environment_

```go
_ = os.Setenv("APP_ENV", "local")
godump.Println(env.IsAppEnv(env.Production, env.Staging))

// #bool false
```

### <a id="isappenvdev"></a>IsAppEnvDev · readonly

IsAppEnvDev checks if APP_ENV is "dev".

```go
_ = os.Setenv("APP_ENV", env.Dev)
godump.Println(env.IsAppEnvDev())

// #bool true
```

### <a id="isappenvlocal"></a>IsAppEnvLocal · readonly

IsAppEnvLocal checks if APP_ENV is "local".

```go
_ = os.Setenv("APP_ENV", env.Local)
godump.Println(env.IsAppEnvLocal())

// #bool true
```

### <a id="isappenvlocalorstaging"></a>IsAppEnvLocalOrStaging · readonly

IsAppEnvLocalOrStaging checks if APP_ENV is either "local" or "staging".

```go
_ = os.Setenv("APP_ENV", env.Local)
godump.Println(env.IsAppEnvLocalOrStaging())

// #bool true
```

### <a id="isappenvproduction"></a>IsAppEnvProduction · readonly

IsAppEnvProduction checks if APP_ENV is "production".

```go
_ = os.Setenv("APP_ENV", env.Production)
godump.Println(env.IsAppEnvProduction())

// #bool true
```

### <a id="isappenvstaging"></a>IsAppEnvStaging · readonly

IsAppEnvStaging checks if APP_ENV is "staging".

```go
_ = os.Setenv("APP_ENV", env.Staging)
godump.Println(env.IsAppEnvStaging())

// #bool true
```

### <a id="isappenvtesting"></a>IsAppEnvTesting · readonly

IsAppEnvTesting reports whether APP_ENV is "testing" or the process looks like `go test`.

_Example: APP_ENV explicitly testing_

```go
_ = os.Setenv("APP_ENV", env.Testing)
godump.Println(env.IsAppEnvTesting())

// #bool true
```

_Example: no test markers_

```go
_ = os.Unsetenv("APP_ENV")
godump.Println(env.IsAppEnvTesting())

// #bool false (outside of test binaries)
```

### <a id="isappenvtestingorlocal"></a>IsAppEnvTestingOrLocal · readonly

IsAppEnvTestingOrLocal checks if APP_ENV is "testing" or "local".

```go
_ = os.Setenv("APP_ENV", env.Testing)
godump.Println(env.IsAppEnvTestingOrLocal())

// #bool true
```

## Other

### <a id="iscontainer"></a>IsContainer

IsContainer detects any container runtime.

### <a id="isdocker"></a>IsDocker

IsDocker reports whether the current process is running in a Docker container.

### <a id="isdockerhost"></a>IsDockerHost

IsDockerHost reports whether this container behaves like a Docker host.

### <a id="isdockerindocker"></a>IsDockerInDocker

IsDockerInDocker reports whether we are inside a Docker-in-Docker environment.

### <a id="isenvloaded"></a>IsEnvLoaded

IsEnvLoaded checks if the environment file has been loaded

### <a id="ishostenvironment"></a>IsHostEnvironment

IsHostEnvironment reports whether the process is running *outside* any
container or orchestrated runtime.

Being a Docker host does NOT count as being in a container.

### <a id="iskubernetes"></a>IsKubernetes

IsKubernetes reports whether running inside Kubernetes.

### <a id="loadenvfileifexists"></a>LoadEnvFileIfExists

LoadEnvFileIfExists loads environment file .env locally
loads .env.testing if invoked from the context of a test file
loads .env.host if invoked from the context of MacOS which references variables to communicate back to the docker network

## Runtime

### <a id="arch"></a>Arch · readonly

Arch returns the CPU architecture the binary is running on.

_Example: print GOARCH_

```go
godump.Println(env.Arch())

// #string "amd64"
// #string "arm64"
```

### <a id="isbsd"></a>IsBSD · readonly

IsBSD reports whether the runtime OS is any BSD variant.

```go
godump.Println(env.IsBSD())

// #bool true  (on BSD variants)
// #bool false (elsewhere)
```

### <a id="iscontaineros"></a>IsContainerOS · readonly

IsContainerOS reports whether this OS is *typically* used as a container base.

```go
godump.Println(env.IsContainerOS())

// #bool true  (on Linux)
// #bool false (on macOS/Windows)
```

### <a id="islinux"></a>IsLinux · readonly

IsLinux reports whether the runtime OS is Linux.

```go
godump.Println(env.IsLinux())

// #bool true  (on Linux)
// #bool false (on other OSes)
```

### <a id="ismac"></a>IsMac · readonly

IsMac reports whether the runtime OS is macOS (Darwin).

```go
godump.Println(env.IsMac())

// #bool true  (on macOS)
// #bool false (elsewhere)
```

### <a id="isunix"></a>IsUnix · readonly

IsUnix reports whether the OS is Unix-like.

```go
godump.Println(env.IsUnix())

// #bool true  (on Unix-like OSes)
// #bool false (e.g., on Windows or Plan 9)
```

### <a id="iswindows"></a>IsWindows · readonly

IsWindows reports whether the runtime OS is Windows.

```go
godump.Println(env.IsWindows())

// #bool true  (on Windows)
// #bool false (elsewhere)
```

### <a id="os"></a>OS · readonly

OS returns the current operating system identifier.

_Example: inspect GOOS_

```go
godump.Println(env.OS())

// #string "linux"   (on Linux)
// #string "darwin"  (on macOS)
// #string "windows" (on Windows)
```

## Typed getters

### <a id="get"></a>Get · readonly

Get returns the environment variable for key or fallback when empty.

_Example: fallback when unset_

```go
os.Unsetenv("DB_HOST")
host := env.Get("DB_HOST", "localhost")
godump.Println(host)

// #string "localhost"
```

_Example: prefer existing value_

```go
_ = os.Setenv("DB_HOST", "db.internal")
host := env.Get("DB_HOST", "localhost")
godump.Println(host)

// #string "db.internal"
```

### <a id="getbool"></a>GetBool · panic

GetBool parses a boolean from an environment variable or fallback string.

_Example: numeric truthy_

```go
_ = os.Setenv("DEBUG", "1")
debug := env.GetBool("DEBUG", "false")
godump.Println(debug)

// #bool true
```

_Example: fallback string_

```go
os.Unsetenv("DEBUG")
debug := env.GetBool("DEBUG", "false")
godump.Println(debug)

// #bool false
```

### <a id="getduration"></a>GetDuration · panic

GetDuration parses a Go duration string (e.g. "5s", "10m", "1h").

_Example: override request timeout_

```go
_ = os.Setenv("HTTP_TIMEOUT", "30s")
timeout := env.GetDuration("HTTP_TIMEOUT", "5s")
godump.Println(timeout)

// #time.Duration 30s
```

_Example: fallback when unset_

```go
os.Unsetenv("HTTP_TIMEOUT")
timeout := env.GetDuration("HTTP_TIMEOUT", "5s")
godump.Println(timeout)

// #time.Duration 5s
```

### <a id="getenum"></a>GetEnum · panic

GetEnum ensures the environment variable's value is in the allowed list.

_Example: accept only staged environments_

```go
_ = os.Setenv("APP_ENV", "prod")
env := env.GetEnum("APP_ENV", "dev", []string{"dev", "staging", "prod"})
godump.Println(env)

// #string "prod"
```

_Example: fallback when unset_

```go
os.Unsetenv("APP_ENV")
env := env.GetEnum("APP_ENV", "dev", []string{"dev", "staging", "prod"})
godump.Println(env)

// #string "dev"
```

### <a id="getfloat"></a>GetFloat · panic

GetFloat parses a float64 from an environment variable or fallback string.

_Example: override threshold_

```go
_ = os.Setenv("THRESHOLD", "0.82")
threshold := env.GetFloat("THRESHOLD", "0.75")
godump.Println(threshold)

// #float64 0.82
```

_Example: fallback with decimal string_

```go
os.Unsetenv("THRESHOLD")
threshold := env.GetFloat("THRESHOLD", "0.75")
godump.Println(threshold)

// #float64 0.75
```

### <a id="getint"></a>GetInt · panic

GetInt parses an int from an environment variable or fallback string.

_Example: fallback used_

```go
os.Unsetenv("PORT")
port := env.GetInt("PORT", "3000")
godump.Println(port)

// #int 3000
```

_Example: env overrides fallback_

```go
_ = os.Setenv("PORT", "8080")
port := env.GetInt("PORT", "3000")
godump.Println(port)

// #int 8080
```

### <a id="getint64"></a>GetInt64 · panic

GetInt64 parses an int64 from an environment variable or fallback string.

_Example: parse large numbers safely_

```go
_ = os.Setenv("MAX_SIZE", "1048576")
size := env.GetInt64("MAX_SIZE", "512")
godump.Println(size)

// #int64 1048576
```

_Example: fallback when unset_

```go
os.Unsetenv("MAX_SIZE")
size := env.GetInt64("MAX_SIZE", "512")
godump.Println(size)

// #int64 512
```

### <a id="getmap"></a>GetMap · readonly

GetMap parses key=value pairs separated by commas into a map.

_Example: parse throttling config_

```go
_ = os.Setenv("LIMITS", "read=10, write=5, burst=20")
limits := env.GetMap("LIMITS", "")
godump.Println(limits)

// #map[string]string [
//  "burst" => "20" #string
//  "read"  => "10" #string
//  "write" => "5" #string
// ]
```

_Example: returns empty map when unset or blank_

```go
os.Unsetenv("LIMITS")
limits := env.GetMap("LIMITS", "")
godump.Println(limits)

// #map[string]string []
```

### <a id="getslice"></a>GetSlice · readonly

GetSlice splits a comma-separated string into a []string with trimming.

_Example: trimmed addresses_

```go
_ = os.Setenv("PEERS", "10.0.0.1, 10.0.0.2")
peers := env.GetSlice("PEERS", "")
godump.Println(peers)

// #[]string [
//  0 => "10.0.0.1" #string
//  1 => "10.0.0.2" #string
// ]
```

_Example: empty becomes empty slice_

```go
os.Unsetenv("PEERS")
peers := env.GetSlice("PEERS", "")
godump.Println(peers)

// #[]string []
```

### <a id="getuint"></a>GetUint · panic

GetUint parses a uint from an environment variable or fallback string.

_Example: defaults to fallback when missing_

```go
os.Unsetenv("WORKERS")
workers := env.GetUint("WORKERS", "4")
godump.Println(workers)

// #uint 4
```

_Example: uses provided unsigned value_

```go
_ = os.Setenv("WORKERS", "16")
workers := env.GetUint("WORKERS", "4")
godump.Println(workers)

// #uint 16
```

### <a id="getuint64"></a>GetUint64 · panic

GetUint64 parses a uint64 from an environment variable or fallback string.

_Example: high range values_

```go
_ = os.Setenv("MAX_ITEMS", "5000")
maxItems := env.GetUint64("MAX_ITEMS", "100")
godump.Println(maxItems)

// #uint64 5000
```

_Example: fallback when unset_

```go
os.Unsetenv("MAX_ITEMS")
maxItems := env.GetUint64("MAX_ITEMS", "100")
godump.Println(maxItems)

// #uint64 100
```

### <a id="mustget"></a>MustGet · panic

MustGet returns the value of key or panics if missing/empty.

_Example: required secret_

```go
_ = os.Setenv("API_SECRET", "s3cr3t")
secret := env.MustGet("API_SECRET")
godump.Println(secret)

// #string "s3cr3t"
```

_Example: panic on missing value_

```go
os.Unsetenv("API_SECRET")
secret := env.MustGet("API_SECRET") // panics: env variable missing: API_SECRET
```

### <a id="mustgetbool"></a>MustGetBool · panic

MustGetBool panics if missing or invalid.

_Example: gate features explicitly_

```go
_ = os.Setenv("FEATURE_ENABLED", "true")
enabled := env.MustGetBool("FEATURE_ENABLED")
godump.Println(enabled)

// #bool true
```

_Example: panic on invalid value_

```go
_ = os.Setenv("FEATURE_ENABLED", "maybe")
_ = env.MustGetBool("FEATURE_ENABLED") // panics when parsing
```

### <a id="mustgetint"></a>MustGetInt · panic

MustGetInt panics if the value is missing or not an int.

_Example: ensure numeric port_

```go
_ = os.Setenv("PORT", "8080")
port := env.MustGetInt("PORT")
godump.Println(port)

// #int 8080
```

_Example: panic on bad value_

```go
_ = os.Setenv("PORT", "not-a-number")
_ = env.MustGetInt("PORT") // panics when parsing
```
<!-- api:embed:end -->

## Philosophy

**env** is part of the **GoForj toolchain** - a collection of focused, composable packages designed to make building Go applications *satisfying*.

No magic. No globals. No surprises.

## License

MIT
