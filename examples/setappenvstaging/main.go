//go:build ignore
// +build ignore

package main

import "github.com/goforj/env/v2"

func main() {
	// SetAppEnvStaging sets APP_ENV to "staging".

	_ = env.SetAppEnvStaging()
	env.Dump(env.GetAppEnv())

	// #string "staging"
}
