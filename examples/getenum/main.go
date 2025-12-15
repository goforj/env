//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env"
	"os"
)

func main() {
	// GetEnum ensures the environment variable's value is in the allowed list.

	// Example: accept only staged environments
	_ = os.Setenv("APP_ENV", "prod")
	appEnv := env.GetEnum("APP_ENV", "dev", []string{"dev", "staging", "prod"})
	env.Dump(appEnv)

	// #string "prod"

	// Example: fallback when unset
	os.Unsetenv("APP_ENV")
	appEnv = env.GetEnum("APP_ENV", "dev", []string{"dev", "staging", "prod"})
	env.Dump(appEnv)

	// #string "dev"
}
