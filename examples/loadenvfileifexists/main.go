//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"os"
	"path/filepath"
)

func main() {
	// LoadEnvFileIfExists loads .env/.env.testing/.env.host when present.

	// Example: test-specific env file
	tmp, _ := os.MkdirTemp("", "envdoc")
	_ = os.WriteFile(filepath.Join(tmp, ".env.testing"), []byte("PORT=9090\nAPP_DEBUG=0"), 0o644)
	_ = os.Chdir(tmp)
	_ = os.Setenv("APP_ENV", env.Testing)

	_ = env.LoadEnvFileIfExists()
	env.Dump(os.Getenv("PORT"))

	// #string "9090"

	// Example: default .env on a host
	_ = os.WriteFile(".env", []byte("SERVICE=api\nAPP_DEBUG=3"), 0o644)
	_ = env.LoadEnvFileIfExists()
	env.Dump(os.Getenv("SERVICE"))

	// #string "api"
}
