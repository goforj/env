//go:build ignore
// +build ignore

package main

import (
	"github.com/goforj/env/v2"
	"os"
)

func main() {
	// GetEnum ensures the environment variable's value is in the allowed list.

	// Example: accept only staged environments
	_ = os.Setenv("APP_ENV", "production")
	appEnv := env.GetEnum("APP_ENV", "local", []string{"local", "staging", "production"})
	env.Dump(appEnv)

	// #string "production"

	// Example: fallback when unset
	os.Unsetenv("APP_ENV")
	appEnv = env.GetEnum("APP_ENV", "local", []string{"local", "staging", "production"})
	env.Dump(appEnv)

	// #string "local"
}
