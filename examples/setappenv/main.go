//go:build ignore
// +build ignore

package main

import "github.com/goforj/env/v2"

func main() {
	// SetAppEnv sets APP_ENV to a supported value.

	// Example: set a supported environment
	_ = env.SetAppEnv(env.Staging)
	env.Dump(env.GetAppEnv())

	// #string "staging"
}
