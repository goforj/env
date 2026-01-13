//go:build ignore
// +build ignore

package main

import "github.com/goforj/env"

func main() {
	// SetAppEnvProduction sets APP_ENV to "production".

	_ = env.SetAppEnvProduction()
	env.Dump(env.GetAppEnv())

	// #string "production"
}
