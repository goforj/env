<p align="center">
  <img src="./docs/images/logo.png" width="500" alt="goforj/env – typed environment variables for Go">
</p>

<p align="center">
  <strong>Typed environment variables for Go with safe defaults.</strong><br>
  Zero-ceremony config access, app environment helpers, and predictable behavior.
</p>

<p align="center">
  <a href="https://pkg.go.dev/github.com/goforj/env">
    <img src="https://pkg.go.dev/badge/github.com/goforj/env.svg" alt="Go Reference">
  </a>
  <a href="LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="MIT License">
  </a>
</p>

---

## Why env?

Accessing environment variables in Go often leads to:

- Repeated parsing logic
- Unsafe string conversions
- Inconsistent defaults
- Scattered app environment checks

**env** solves this by providing **typed accessors with fallbacks**, so configuration stays boring and predictable.

---

## Features

- Strongly typed getters (`int`, `bool`, `duration`, slices, maps)
- Safe fallbacks (never panic, never empty by accident)
- App environment helpers (`dev`, `local`, `prod`)
- Zero dependencies
- Framework-agnostic

---

## Installation

```bash
go get github.com/goforj/env
````

---

## Usage

```go
import "github.com/goforj/env"

port := env.GetInt("PORT", 8080)
debug := env.GetBool("DEBUG", false)
timeout := env.GetDuration("REQUEST_TIMEOUT", time.Second*5)
```

### Application environment

```go
if env.IsAppEnvDev() {
    // dev-only behavior
}
```

---

## API Overview

```go
env.Get(key, fallback string) string
env.GetBool(key string, fallback bool) bool
env.GetInt(key string, fallback int) int
env.GetInt64(key string, fallback int64) int64
env.GetUint(key string, fallback uint) uint
env.GetUint64(key string, fallback uint64) uint64
env.GetFloat(key string, fallback float64) float64
env.GetDuration(key string, fallback time.Duration) time.Duration
env.GetSlice(key string, fallback []string) []string
env.GetMap(key string, fallback map[string]string) map[string]string
env.GetEnum(key string, fallback string, allowed []string) string

env.GetAppEnv() string
env.IsAppEnv(...string) bool
env.IsAppEnvDev() bool
env.IsAppEnvLocal() bool
env.IsAppEnvProd() bool
```

---

## Philosophy

**env** is part of the **GoForj toolchain** — a collection of focused, composable packages designed to make building Go applications *satisfying*.

No magic. No globals. No surprises.

---

## License

MIT
